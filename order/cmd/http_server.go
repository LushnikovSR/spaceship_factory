package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	customMiddleware "github.com/LushnikovSR/spaceship_factory/internal/middleware"
	apiV1 "github.com/LushnikovSR/spaceship_factory/order/internal/api/order/v1"
	inventoryClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/payment/v1"
	config "github.com/LushnikovSR/spaceship_factory/order/internal/config"
	repository "github.com/LushnikovSR/spaceship_factory/order/internal/repository/order"
	service "github.com/LushnikovSR/spaceship_factory/order/internal/service/order"
	orderV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

const (
	contextTimeout = 10 * time.Second
	configPath     = "./deploy/compose/order/.env"
)

func main() {
	err := godotenv.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}
	if err := run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

func run() error {
	// Контекст для подключения к PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	// Создаём соединение с базой данных
	pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer pool.Close()

	// Проверяем что соединение с базой данных установлено
	err = pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database is unavailable: %w", err)
	}

	// Создаём новое хранилище для данных о заказах
	repo, err := repository.NewRepository(pool)
	if err != nil {
		return fmt.Errorf("failed to init repository: %w", err)
	}

	// gRPC подключение к InventoryService
	connInv, err := grpc.NewClient(config.AppConfig().InventoryGRPC.Address(),
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
	connPay, err := grpc.NewClient(config.AppConfig().PaymentGRPC.Address(),
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
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server is starting on port: %s\n", config.AppConfig().OrderHTTP.Address())
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
