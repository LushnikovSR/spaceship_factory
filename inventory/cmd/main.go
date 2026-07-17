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
	"time"

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

func run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(panicRecoveryInterceptor()))

	// Контекст для подключения к MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := repository.ConnectMongo(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer func() {
		if cerr := client.Disconnect(context.Background()); cerr != nil {
			slog.Warn("MongoDB disconnect error", "error", cerr)
		}
	}()

	// Пинг с отдельным коротким таймаутом
	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db := client.Database("inventory")
	repo := repository.NewRepository(db)
	repo.Init(ctx)

	svc := service.NewService(repo)
	api := apiV1.NewAPI(svc)
	inventory_v1.RegisterInventoryServiceServer(s, api)
	reflection.Register(s)

	go func() {
		slog.Info("gRPC server listening", "port", grpcPort)
		if err := s.Serve(lis); err != nil {
			slog.Error("failed to serve", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down gRPC server...")
	s.GracefulStop()
	slog.Info("Server stopped")
	return nil
}

func main() {
	if err := run(); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}
