package order

import (
	grpc "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc"
	repository "github.com/LushnikovSR/spaceship_factory/order/internal/repository"
	def "github.com/LushnikovSR/spaceship_factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient

	orderRepository repository.OrderRepository
}

func NewService(orderRepository repository.OrderRepository,
	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *service {
	return &service{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderRepository: orderRepository,
	}
}
