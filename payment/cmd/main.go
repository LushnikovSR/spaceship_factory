package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	apiV1 "github.com/LushnikovSR/spaceship_factory/payment/internal/api/payment/v1"
	service "github.com/LushnikovSR/spaceship_factory/payment/internal/service/payment"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

const (
	grpcPort = 50052
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	defer func() {
		err := lis.Close()
		if err != nil {
			slog.Error("failed to close listener", "error", err)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := service.NewService()
	api := apiV1.NewAPI(service)

	payment_v1.RegisterPaymentServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		slog.Info("🚀 gRPC server listening", "port", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			slog.Error("failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	slog.Info("✅ Server stopped")
}
