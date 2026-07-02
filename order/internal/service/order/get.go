package order

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

// GetOrder implements getOrder operation.
//
// Get order by uuid.
//
// GET /orders/{order_uuid}
func (s *service) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order := s.orderRepository.GetOrder(orderUUID)
	if order == nil {
		return model.Order{}, &model.NotFoundError{
			BaseError: model.BaseError{
				Code:    404,
				Message: "Order for uuid '" + orderUUID + "' not found",
			},
		}
	}

	return *order, nil
}
