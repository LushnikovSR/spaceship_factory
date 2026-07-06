package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiV1 "github.com/LushnikovSR/spaceship_factory/order/internal/api/order/v1"
	inventoryClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/payment/v1"
	repository "github.com/LushnikovSR/spaceship_factory/order/internal/repository/order"
	service "github.com/LushnikovSR/spaceship_factory/order/internal/service/order"

	customMiddleware "github.com/LushnikovSR/spaceship_factory/internal/middleware"
	orderV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort             = "8080"
	readHeaderTimeout    = 5 * time.Second
	contextTimeout       = 10 * time.Second
	inventoryServicePort = "50051"
	paymentServicePort   = "50052"
)

func main() {
	//Создаём новое хранилище для данных о заказах
	repo := repository.NewRepository()

	// gRPC подключение к InventoryService
	connInv, err := grpc.NewClient(fmt.Sprintf("localhost:%s", inventoryServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to InventoryService: %v", err)
	}
	defer func() {
		if cerr := connInv.Close(); cerr != nil {
			log.Printf("failed to close InventoryService connection: %v", cerr)
		}
	}()
	invClient := inventory_v1.NewInventoryServiceClient(connInv)

	//gRPC подключение к PaymentService
	connPay, err := grpc.NewClient(fmt.Sprintf("localhost:%s", paymentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to PaymentService: %v", err)
	}
	defer func() {
		if cerr := connPay.Close(); cerr != nil {
			log.Printf("failed to close PaymentService connection: %v", cerr)
		}
	}()
	payClient := payment_v1.NewPaymentServiceClient(connPay)

	//Создаём клиентов для inventory, payment сервисов
	inventoryClient := inventoryClient.NewClient(invClient)
	paymentClient := paymentClient.NewClient(payClient)

	//Создаём order сервис с api
	service := service.NewService(repo, inventoryClient, paymentClient)
	api := apiV1.NewAPI(service)

	//Создаём ОpenAPI сервер
	orderServer, err := orderV1.NewServer(api, orderV1.WithPathPrefix("/api/v1"))
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	//Добавляем middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(customMiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	//Запускаем http-server
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, //для защиты от slowloris атак
	}

	//Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server is starting on port: %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Starting server is failed: %s\n", err)
		}
	}()

	//Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server stopping ...")

	//Создаём контекст с таймутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Error during server closing: %v\n", err)
	}
	log.Printf("Server is stoped correctly")
}
