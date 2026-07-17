package order

import (
	"context"
	"fmt"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
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
func (s *service) PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID string) (string, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		return "", &model.InternalServerError{
			BaseError: model.BaseError{
				Code:    500,
				Message: fmt.Errorf("order for uuid %v not found: %w", orderUUID, err).Error(),
			},
		}
	}
	if order == nil {
		return "", &model.NotFoundError{
			BaseError: model.BaseError{
				Code:    404,
				Message: "Order for uuid '" + orderUUID + "' not found",
			},
		}
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, order.OrderUUID, order.UserUUID, string(paymentMethod))
	if err != nil {
		return "", &model.InternalServerError{
			BaseError: model.BaseError{
				Code:    500,
				Message: "Payment service error: " + err.Error(),
			},
		}
	}

	order.SetTransactionUUID(model.OptNilString{Value: transactionUUID, Set: true})
	order.SetPaymentMethod(&model.NilOrderDtoPaymentMethod{Value: model.OrderDtoPaymentMethod(paymentMethod)})
	order.SetStatus(model.OrderDtoStatusPAID)

	err = s.orderRepository.UpdateOrder(ctx, order)
	if err != nil {
		return "", &model.InternalServerError{
			BaseError: model.BaseError{
				Code:    500,
				Message: fmt.Errorf("The order status wasn`t update to PAID: %w", err).Error(),
			},
		}
	}

	return transactionUUID, nil
}
