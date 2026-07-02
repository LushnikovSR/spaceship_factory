package order

import (
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/order/internal/repository/model"
)

func OrderModelToRepoModel(order *model.Order) *repoModel.Order {
	if order == nil {
		return nil
	}
	var repoPayment *repoModel.NilOrderDtoPaymentMethod
	if order.PaymentMethod != nil {
		// Конвертируем enum во внутреннем поле
		repoPayment = &repoModel.NilOrderDtoPaymentMethod{
			Value: repoModel.OrderDtoPaymentMethod(order.PaymentMethod.Value),
			Null:  order.PaymentMethod.Null,
		}
	}

	return &repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: repoModel.OptNilString(order.TransactionUUID),
		PaymentMethod:   repoPayment,
		Status:          repoModel.OrderDtoStatus(order.Status),
	}
}

func OrderRepoModelToModel(order *repoModel.Order) *model.Order {
	if order == nil {
		return nil
	}
	var payment *model.NilOrderDtoPaymentMethod
	if order.PaymentMethod != nil {
		// Конвертируем enum во внутреннем поле
		payment = &model.NilOrderDtoPaymentMethod{
			Value: model.OrderDtoPaymentMethod(order.PaymentMethod.Value),
			Null:  order.PaymentMethod.Null,
		}
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
