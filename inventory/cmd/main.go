package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	apiV1 "github.com/LushnikovSR/spaceship_factory/inventory/internal/api/inventory/v1"
	repository "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/part"
	service "github.com/LushnikovSR/spaceship_factory/inventory/internal/service/part"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

const (
	grpcPort = 50051
)

func panicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("PANIC",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

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
	s := grpc.NewServer(
		grpc.UnaryInterceptor(panicRecoveryInterceptor()),
	)

	// Регистрируем Inventory сервис
	repo := repository.NewRepository()

	repo.Init()

	service := service.NewService(repo)
	api := apiV1.NewAPI(service)

	inventory_v1.RegisterInventoryServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		fmt.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			slog.Error("failed to serve", "error", err)
			os.Exit(1)
		}
	}()

	// Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	slog.Info("✅ Server stopped")
}
