package order

import (
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	order_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
)

func OrderModelToAPI(order *model.Order) *order_v1.OrderDto {
	payment := &order_v1.NilOrderDtoPaymentMethod{
		Value: order_v1.OrderDtoPaymentMethod(order.PaymentMethod.Value),
		Null:  order.PaymentMethod.Null,
	}

	return &order_v1.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order_v1.OptNilString(order.TransactionUUID),
		PaymentMethod:   payment,
		Status:          order_v1.OrderDtoStatus(order.Status),
	}
}

func OrderAPIToModel(order *order_v1.OrderDto) *model.Order {
	payment := &model.NilOrderDtoPaymentMethod{
		Value: model.OrderDtoPaymentMethod(order.PaymentMethod.Value),
		Null:  order.PaymentMethod.Null,
	}

	return &model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: model.OptNilString(order.TransactionUUID),
		PaymentMethod:   payment,
		Status:          model.OrderDtoStatus(order.Status),
	}
}
