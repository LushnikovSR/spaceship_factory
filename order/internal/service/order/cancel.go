package order

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

// CancelOrder implements cancelOrder operation.
//
// Checks the order status. If `PENDING_PAYMENT`, changes the status to `CANCELLED`. If `PAID`,
// returns a 409 error.
//
// POST /orders/{order_uuid}/cancel
func (s *service) CancelOrder(ctx context.Context, uuid string) error {
	order := s.orderRepository.GetOrder(uuid)
	if order == nil {
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
				Message: "The order has been paid. Cancellation is not possible.",
			},
		}
	}

	if order.Status == model.OrderDtoStatusCANCELLED {
		return &model.BadRequestError{
			BaseError: model.BaseError{
				Code:    400,
				Message: "The order has already been cancelled. Cancellation is not possible again.",
			},
		}
	}

	order.SetStatus(model.OrderDtoStatusCANCELLED)
	s.orderRepository.UpdateOrder(order)

	return nil
}
