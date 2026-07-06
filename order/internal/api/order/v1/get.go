package order

import (
	"context"

	converter "github.com/LushnikovSR/spaceship_factory/order/internal/converter"
	order_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
)

// GetOrder implements getOrder operation.
//
// Get order by uuid.
//
// GET /orders/{order_uuid}
func (a *api) GetOrder(ctx context.Context, params order_v1.GetOrderParams) (order_v1.GetOrderRes, error) {
	order, err := a.orderService.GetOrder(ctx, params.OrderUUID)
	if err != nil {
		return &order_v1.OrderDto{}, err
	}
	return converter.OrderModelToAPI(&order), nil
}
