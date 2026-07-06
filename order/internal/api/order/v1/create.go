package order

import (
	"context"

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
		return &order_v1.CreateOrderResponse{}, err
	}
	return &order_v1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: total_price,
	}, nil
}
