package order

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/order/internal/repository/model"
)

func OrderModelToRepoModel(order *model.Order) *repoModel.Order {
	if order == nil {
		return nil
	}

	var transactionUUID *string
	if order.TransactionUUID.Set && !order.TransactionUUID.Null {
		transactionUUID = &order.TransactionUUID.Value
	}

	var paymentMethod string
	if !order.PaymentMethod.Null {
		paymentMethod = string(order.PaymentMethod.Value)
	}

	return &repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          string(order.Status),
	}
}

func OrderRepoModelToModel(order *repoModel.Order) *model.Order {
	if order == nil {
		return nil
	}

	var transactionID model.OptNilString
	if order.TransactionUUID != nil {
		transactionID.SetTo(*order.TransactionUUID)
	} else {
		transactionID.SetToNull()
	}

	paymentMethod := &model.NilOrderDtoPaymentMethod{}
	if order.PaymentMethod != nil {
		paymentMethod.Value = model.OrderDtoPaymentMethod(*order.PaymentMethod)
	} else {
		paymentMethod.Null = true
	}

	return &model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionID,
		PaymentMethod:   paymentMethod,
		Status:          model.OrderDtoStatus(order.Status),
	}
}

func UuidToRepo(uuid string) (primitive.ObjectID, error) {
	repoUUID, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		return primitive.ObjectID{}, fmt.Errorf("failed to convert to primitive.ObjectID: %w", err)
	}
	return repoUUID, nil
}
