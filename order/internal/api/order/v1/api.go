package order

import (
	service "github.com/LushnikovSR/spaceship_factory/order/internal/service"
	order_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
)

var _ order_v1.Handler = (*api)(nil)

type api struct {
	order_v1.UnimplementedHandler

	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
