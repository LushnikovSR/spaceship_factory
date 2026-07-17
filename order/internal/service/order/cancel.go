package order

import (
	"context"
	"fmt"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

// CancelOrder implements cancelOrder operation.
//
// Checks the order status. If `PENDING_PAYMENT`, changes the status to `CANCELLED`. If `PAID`,
// returns a 409 error.
//
// POST /orders/{order_uuid}/cancel
func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	order, err := s.orderRepository.GetOrder(ctx, uuid)
	if err != nil {
		return err
	}
	if order.OrderUUID == "" {
		return &model.NotFoundError{
			BaseError: model.BaseError{
				Code:    404,
				Message: "Order for uuid '" + uuid + "' not found",
			},
		}
	}

	if order.Status == model.OrderDtoStatusPAID {
		return &model.ConflictError{
			BaseError: model.BaseError{
				Code:    409,
				Message: "The order '" + uuid + "' has been paid. Cancellation is not possible.",
			},
		}
	}

	if order.Status == model.OrderDtoStatusCANCELLED {
		return &model.BadRequestError{
			BaseError: model.BaseError{
				Code:    400,
				Message: "The order '" + uuid + "' has already been cancelled. Cancellation is not possible again.",
			},
		}
	}

	order.SetStatus(model.OrderDtoStatusCANCELLED)
	err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return &model.InternalServerError{
			BaseError: model.BaseError{
				Code:    500,
				Message: fmt.Errorf("The order status wasn`t update to CANCELLED: %w", err).Error(),
			},
		}
	}

	return nil
}
