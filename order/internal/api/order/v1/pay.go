package order

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	order_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
)

// PayOrder implements payOrder operation.
//
// Находит заказ по `order_uuid`. Если не существует —
// возвращает 404 Not Found. Вызывает `PaymentService.PayOrder`, передаёт
// `user_uuid`, `order_uuid` и `payment_method`. Получает`transaction_uuid`.
// Обновляет заказ: статус → `PAID`, сохраняет `transaction_uuid`,
// `payment_method`.
//
// POST /orders/{order_uuid}/pay
func (a *api) PayOrder(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderParams) (order_v1.PayOrderRes, error) {

	transaction_UUID, err := a.orderService.PayOrder(ctx, model.PaymentMethod(req.PaymentMethod), params.OrderUUID)
	if err != nil {
		return &order_v1.PayOrderResponse{}, err
	}
	return &order_v1.PayOrderResponse{
		TransactionUUID: transaction_UUID,
	}, nil
}
