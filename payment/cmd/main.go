package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	grpcPort = 50052
)

//PaymentService реализует gRPC сервис для проведения оплаты заказов
type PaymentService struct {
	payment_v1.UnimplementedPaymentServiceServer
}

//PayOrder проверяет Request на наличие входных данных и возвращает Response с uuid транзакции.
func (s *PaymentService) PayOrder(_ context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	// 1. Валидация обязательных полей
	if req.OrderUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "order_uuid is required")
	}
	if req.UserUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	// 2. Проверка, что метод оплаты указан (не нулевое "неопределённое" значение)
	if req.PaymentMethod == payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED { // предполагаем, что 0 = UNSPECIFIED
		return nil, status.Error(codes.InvalidArgument, "payment_method must be specified")
	}

	transactionUuid := uuid.NewString()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUuid)

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
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
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &PaymentService{}

	payment_v1.RegisterPaymentServiceServer(s, service)

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
