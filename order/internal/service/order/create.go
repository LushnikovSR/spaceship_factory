package order

import (
	"context"
	"fmt"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

// CreateOrder implements createOrder operation.
//
// Получает детали через `InventoryService.ListParts`. Проверяет, что
// все детали существуют. Если хотя бы одной нет —
// возвращает ошибку. Считает `total_price`. Генерирует `order_uuid`.
// Сохраняет заказ со статусом `PENDING_PAYMENT`.
//
// POST /orders
func (s *service) CreateOrder(ctx context.Context, userUUID string, partUuids []string) (orderUUID string, totalPrice float64, err error) {
	// Получаем детали из InventoryService
	parts, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		Uuids: partUuids,
	})
	if err != nil {
		// Ошибка связи с InventoryService
		return "", 0, &model.InternalServerError{
			BaseError: model.BaseError{
				Code:    500,
				Message: "Failed to fetch parts from inventory",
			},
		}
	}

	// Проверяем, что все запрошенные детали существуют
	foundUuids := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		foundUuids[p.Uuid] = struct{}{}
	}
	for _, uid := range partUuids {
		if _, ok := foundUuids[uid]; !ok {
			return "", 0, &model.NotFoundError{
				BaseError: model.BaseError{
					Code:    404,
					Message: "Part with UUID " + uid + " not found",
				},
			}
		}
	}

	// Считаем total_price
	total_price := float64(0.0)
	for _, part := range parts {
		total_price += part.Price
	}

	var transactionUUID model.OptNilString
	transactionUUID.SetToNull()

	var paymentMethod model.NilOrderDtoPaymentMethod
	paymentMethod.SetToNull()

	order := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUuids,
		TotalPrice:      total_price,
		TransactionUUID: transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          model.OrderDtoStatusPENDINGPAYMENT,
	}

	// Сохраняем данные о заказе в аргументе order
	err = s.orderRepository.CreateOrder(ctx, &order)
	if err != nil {
		return "", 0, &model.ConflictError{
			BaseError: model.BaseError{
				Code:    409,
				Message: fmt.Errorf("Failed to create order in the database: %w", err).Error(),
			},
		}
	}

	return order.OrderUUID, order.TotalPrice, nil
}
