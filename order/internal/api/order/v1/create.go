package order

import (
	"context"
	"errors"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	order_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
)

// CreateOrder implements createOrder operation.
//
// Получает детали через `InventoryService.ListParts`. Проверяет, что
// все детали существуют. Если хотя бы одной нет —
// возвращает ошибку. Считает `total_price`. Генерирует `order_uuid`.
//
// Сохраняет заказ со статусом `PENDING_PAYMENT`.
//
// POST /orders
func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	orderUUID, total_price, err := a.orderService.CreateOrder(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		// Проверяем известные типы ошибок с учётом возможных обёрток
		var notFoundErr *model.NotFoundError
		if errors.As(err, &notFoundErr) {
			return &order_v1.NotFoundError{
				Code:    notFoundErr.Code,
				Message: notFoundErr.Message,
			}, nil
		}

		var conflictErr *model.ConflictError
		if errors.As(err, &conflictErr) {
			return &order_v1.ConflictError{
				Code:    conflictErr.Code,
				Message: conflictErr.Message,
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

	return &order_v1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: total_price,
	}, nil
}
