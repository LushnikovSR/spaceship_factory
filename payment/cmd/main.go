package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	apiV1 "github.com/LushnikovSR/spaceship_factory/payment/internal/api/payment/v1"
	service "github.com/LushnikovSR/spaceship_factory/payment/internal/service/payment"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	grpcPort = 50052
)

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
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := service.NewService()
	api := apiV1.NewAPI(service)

	payment_v1.RegisterPaymentServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
