package order

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	order_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
)

// CreateOrder implements createOrder operation.
//
// Получает детали через `InventoryService.ListParts`. Проверяет, что
// все детали существуют. Если хотя бы одной нет —
// возвращает ошибку. Считает `total_price`. Генерирует `order_uuid`.
//
//	Сохраняет заказ со статусом `PENDING_PAYMENT`.
//
// POST /orders
func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	orderUUID, total_price, err := a.orderService.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		// Конвертируем внутреннюю ошибку в соответствующий тип ответа
		switch e := err.(type) {
		case *model.NotFoundError:
			return &order_v1.NotFoundError{
				Code:    e.Code,
				Message: e.Message,
			}, nil
		case *model.ConflictError:
			return &order_v1.ConflictError{
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
	return &order_v1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: total_price,
	}, nil
}
