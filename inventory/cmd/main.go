package inventory

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	apiV1 "github.com/LushnikovSR/spaceship_factory/inventory/internal/api/inventory/v1"
	repository "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/part"
	service "github.com/LushnikovSR/spaceship_factory/inventory/internal/service/part"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	grpcPort = 50051
)

func panicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC in %s: %v\n%s", info.FullMethod, r, string(debug.Stack()))
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		err := lis.Close()
		if err != nil {
			fmt.Printf("failed to close listener: %v\n", err)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer(
		grpc.UnaryInterceptor(panicRecoveryInterceptor()),
	)

	//Регистрируем Inventory сервис
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
			fmt.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	//Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
