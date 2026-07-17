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
	order := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1", "part2"},
		TotalPrice:      1234.56,
		TransactionUUID: model.NewOptNilString("txn-uuid"),
		PaymentMethod:   nil,
		Status:          model.OrderDtoStatusPENDINGPAYMENT,
	}

	transactionID := "txn-uuid"

	expected := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1", "part2"},
		TotalPrice:      1234.56,
		TransactionUUID: &transactionID,
		PaymentMethod:   nil,
		Status:          "PENDING_PAYMENT",
	}

	result := OrderModelToRepoModel(order)
	assert.Equal(t, expected, result)
}

func TestOrderModelToRepoModel_WithPaymentMethod(t *testing.T) {
	pm := model.NewNilOrderDtoPaymentMethod(model.OrderDtoPaymentMethodCARD)
	order := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1"},
		TotalPrice:      99.99,
		TransactionUUID: model.OptNilString{Set: false}, // не задано
		PaymentMethod:   &pm,
		Status:          model.OrderDtoStatusPAID,
	}

	expectedPM := "CARD"
	expected := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"part1"},
		TotalPrice:      99.99,
		TransactionUUID: nil,
		PaymentMethod:   &expectedPM,
		Status:          "PAID",
	}

	result := OrderModelToRepoModel(order)
	assert.Equal(t, expected, result)
}

func TestOrderModelToRepoModel_PaymentMethodNullTrue(t *testing.T) {
	pm := model.NilOrderDtoPaymentMethod{Null: true}
	order := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       nil,
		TotalPrice:      0,
		TransactionUUID: model.OptNilString{},
		PaymentMethod:   &pm,
		Status:          model.OrderDtoStatusCANCELLED,
	}

	expected := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       nil,
		TotalPrice:      0,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          "CANCELLED",
	}

	result := OrderModelToRepoModel(order)
	assert.Equal(t, expected, result)
}

func TestOrderRepoModelToModel_NilInput(t *testing.T) {
	assert.Nil(t, OrderRepoModelToModel(nil))
}

func TestOrderRepoModelToModel_NoPaymentMethod(t *testing.T) {
	transactionID := "txn-uuid"
	repoOrder := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partA", "partB"},
		TotalPrice:      42.0,
		TransactionUUID: &transactionID,
		PaymentMethod:   nil,
		Status:          "PAID",
	}

	expected := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partA", "partB"},
		TotalPrice:      42.0,
		TransactionUUID: model.NewOptNilString("txn-uuid"),
		PaymentMethod:   nil,
		Status:          model.OrderDtoStatusPAID,
	}

	result := OrderRepoModelToModel(repoOrder)
	assert.Equal(t, expected, result)
}

func TestOrderRepoModelToModel_WithPaymentMethod(t *testing.T) {
	paymentMethod := "SPB"
	repoOrder := &repoModel.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partX"},
		TotalPrice:      500.0,
		TransactionUUID: nil,
		PaymentMethod:   &paymentMethod,
		Status:          "CANCELLED",
	}

	expectedPM := model.NewNilOrderDtoPaymentMethod(model.OrderDtoPaymentMethodSBP)
	expected := &model.Order{
		OrderUUID:       "order-uuid",
		UserUUID:        "user-uuid",
		PartUuids:       []string{"partX"},
		TotalPrice:      500.0,
		TransactionUUID: model.OptNilString{Set: true, Null: true},
		PaymentMethod:   &expectedPM,
		Status:          model.OrderDtoStatusCANCELLED,
	}

	result := OrderRepoModelToModel(repoOrder)
	assert.Equal(t, expected, result)
}
