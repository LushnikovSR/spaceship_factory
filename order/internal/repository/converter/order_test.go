package order

import (
	"testing"

	"github.com/stretchr/testify/assert"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/order/internal/repository/model"
)

func TestOrderModelToRepoModel_NilInput(t *testing.T) {
	assert.Nil(t, OrderModelToRepoModel(nil))
}

func TestOrderModelToRepoModel_NoPaymentMethod(t *testing.T) {
	paymentMethod := &model.NilOrderDtoPaymentMethod{}
	paymentMethod.SetToNull()

	order := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1", "part2"},
		TotalPrice:      1234.56,
		TransactionUUID: model.NewOptNilString("txn-uuid"),
		PaymentMethod:   paymentMethod,
		Status:          model.OrderDtoStatusPENDINGPAYMENT,
	}

	transactionID := "txn-uuid"
	var expectedPaymentMethod string

	expected := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1", "part2"},
		TotalPrice:      1234.56,
		TransactionUUID: &transactionID,
		PaymentMethod:   &expectedPaymentMethod,
		Status:          "PENDING_PAYMENT",
	}

	result := OrderModelToRepoModel(order)
	assert.Equal(t, *expected, *result)
}

func TestOrderModelToRepoModel_WithPaymentMethod(t *testing.T) {
	pm := model.NewNilOrderDtoPaymentMethod(model.OrderDtoPaymentMethodCARD)
	transactionID := model.OptNilString{}
	transactionID.SetToNull()

	order := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1"},
		TotalPrice:      99.99,
		TransactionUUID: transactionID, // не задано
		PaymentMethod:   &pm,
		Status:          model.OrderDtoStatusPAID,
	}

	expectedPM := "CARD"
	var expectedTransactionID *string

	expected := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1"},
		TotalPrice:      99.99,
		TransactionUUID: expectedTransactionID,
		PaymentMethod:   &expectedPM,
		Status:          "PAID",
	}

	result := OrderModelToRepoModel(order)
	assert.Equal(t, expected, result)
}

func TestOrderRepoModelToModel_NilInput(t *testing.T) {
	assert.Nil(t, OrderRepoModelToModel(nil))
}

func TestOrderRepoModelToModel_NoPaymentMethod(t *testing.T) {
	transactionID := "txn-uuid"
	var paymentMethod *string

	repoOrder := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partA", "partB"},
		TotalPrice:      42.0,
		TransactionUUID: &transactionID,
		PaymentMethod:   paymentMethod,
		Status:          "PAID",
	}

	expectedPaymentMethod := &model.NilOrderDtoPaymentMethod{}
	expectedPaymentMethod.SetToNull()

	expected := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partA", "partB"},
		TotalPrice:      42.0,
		TransactionUUID: model.NewOptNilString("txn-uuid"),
		PaymentMethod:   expectedPaymentMethod,
		Status:          model.OrderDtoStatusPAID,
	}

	result := OrderRepoModelToModel(repoOrder)
	assert.Equal(t, expected, result)
}

func TestOrderRepoModelToModel_WithPaymentMethod(t *testing.T) {
	paymentMethod := "SBP"
	repoOrder := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partX"},
		TotalPrice:      500.0,
		TransactionUUID: nil,
		PaymentMethod:   &paymentMethod,
		Status:          "CANCELLED",
	}

	expectedPM := &model.NilOrderDtoPaymentMethod{}
	expectedPM.SetTo(model.OrderDtoPaymentMethodSBP)

	expected := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partX"},
		TotalPrice:      500.0,
		TransactionUUID: model.OptNilString{Set: true, Null: true},
		PaymentMethod:   expectedPM,
		Status:          model.OrderDtoStatusCANCELLED,
	}

	result := OrderRepoModelToModel(repoOrder)
	assert.Equal(t, expected, result)
}
