package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/LushnikovSR/spaceship_factory/internal/middleware"
	apiV1 "github.com/LushnikovSR/spaceship_factory/order/internal/api/order/v1"
	inventoryClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/payment/v1"
	repository "github.com/LushnikovSR/spaceship_factory/order/internal/repository/order"
	service "github.com/LushnikovSR/spaceship_factory/order/internal/service/order"
	orderV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

const (
	httpPort             = "8080"
	readHeaderTimeout    = 5 * time.Second
	contextTimeout       = 10 * time.Second
	inventoryServicePort = "50051"
	paymentServicePort   = "50052"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := repository.ConnectPostgres(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	// Создаём новое хранилище для данных о заказах
	repo, err := repository.NewRepository(pool)
	if err != nil {
		return fmt.Errorf("failed to init repository: %w", err)
	}

	// gRPC подключение к InventoryService
	connInv, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%s", inventoryServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to InventoryService: %w", err)
	}
	defer func() {
		if cerr := connInv.Close(); cerr != nil {
			slog.Error("failed to close InventoryService connection", "error", cerr)
		}
	}()
	invClient := inventory_v1.NewInventoryServiceClient(connInv)

	// gRPC подключение к PaymentService
	connPay, err := grpc.NewClient(fmt.Sprintf("127.0.0.1:%s", paymentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to PaymentService: %w", err)
	}
	defer func() {
		if cerr := connPay.Close(); cerr != nil {
			slog.Error("failed to close PaymentService connection", "error", cerr)
		}
	}()
	payClient := payment_v1.NewPaymentServiceClient(connPay)

	// Создаём клиентов для inventory, payment сервисов
	inventoryClient := inventoryClient.NewClient(invClient)
	paymentClient := paymentClient.NewClient(payClient)

	// Создаём order сервис с api
	svc := service.NewService(repo, inventoryClient, paymentClient)
	api := apiV1.NewAPI(svc)

	// Создаём OpenAPI сервер
	orderServer, err := orderV1.NewServer(api, orderV1.WithPathPrefix("/api/v1"))
	if err != nil {
		return fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(customMiddleware.RequestLogger)
	r.Mount("/", orderServer)

	// Запускаем http-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server is starting on port: %s\n", httpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("🛑 Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("error during server closing: %w", err)
	}
	slog.Info("✅ Server stopped")
	return nil
}
