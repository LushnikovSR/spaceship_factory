package order

import (
	"context"

	order_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
)

// CancelOrder implements cancelOrder operation.
//
// Checks the order status. If `PENDING_PAYMENT`, changes the status to `CANCELLED`. If `PAID`,
// returns a 409 error.
//
// POST /orders/{order_uuid}/cancel
func (a *api) CancelOrder(ctx context.Context, params order_v1.CancelOrderParams) (order_v1.CancelOrderRes, error) {
	err := a.orderService.CancelOrder(ctx, params.OrderUUID)
	if err != nil {
		return &order_v1.CancelOrderNoContent{}, err
	}
	return &order_v1.CancelOrderNoContent{}, nil
}
