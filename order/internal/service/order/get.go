package order

import (
	"context"
	"fmt"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

// GetOrder implements getOrder operation.
//
// Get order by uuid.
//
// GET /orders/{order_uuid}
func (s *service) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return model.Order{}, &model.InternalServerError{
			BaseError: model.BaseError{
				Code:    500,
				Message: fmt.Errorf("order for uuid %v not found: %w", orderUUID, err).Error(),
			},
		}
	}
	if order == nil {
		return model.Order{}, &model.NotFoundError{
			BaseError: model.BaseError{
				Code:    404,
				Message: "order for uuid '" + orderUUID + "' not found",
			},
		}
	}

	return *order, nil
}
