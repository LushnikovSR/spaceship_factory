package order

import (
	"context"

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
		// Конвертируем внутреннюю ошибку в соответствующий тип ответа
		switch e := err.(type) {
		case *model.NotFoundError:
			return &order_v1.NotFoundError{
				Code:    e.Code,
				Message: e.Message,
			}, nil
		case *model.InternalServerError:
			return &order_v1.InternalServerError{
				Code:    e.Code,
				Message: e.Message,
			}, nil
		default:
			return &order_v1.InternalServerError{
				Code:    500,
				Message: err.Error(),
			}, nil
		}
	}
	return converter.OrderModelToAPI(&order), nil
}
