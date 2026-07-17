package order

import (
	"context"
	"errors"

	converter "github.com/LushnikovSR/spaceship_factory/order/internal/converter"
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
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
		// Проверяем известные типы ошибок с учётом возможных обёрток
		var notFoundErr *model.NotFoundError
		if errors.As(err, &notFoundErr) {
			return &order_v1.NotFoundError{
				Code:    notFoundErr.Code,
				Message: notFoundErr.Message,
			}, nil
		}

		var internalErr *model.InternalServerError
		if errors.As(err, &internalErr) {
			return &order_v1.InternalServerError{
				Code:    internalErr.Code,
				Message: internalErr.Message,
			}, nil
		}

		// Неизвестная ошибка
		return &order_v1.InternalServerError{
			Code:    500,
			Message: err.Error(),
		}, nil
	}

	return converter.OrderModelToAPI(&order), nil
}
