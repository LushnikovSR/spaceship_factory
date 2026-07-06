package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	var (
		expectedOrderUUID = gofakeit.UUID()
		userUUID          = gofakeit.UUID()
		total             = 1800.0
		partUuids         = []string{
			"11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222",
		}
		expectedTransactionUUID = gofakeit.UUID()

		expectedPaymentMethod = model.PaymentMethodCARD

		expectedOrder = &model.Order{
			OrderUUID:       expectedOrderUUID,
			UserUUID:        userUUID,
			PartUuids:       partUuids,
			TotalPrice:      total,
			TransactionUUID: model.OptNilString{},
			PaymentMethod:   &model.NilOrderDtoPaymentMethod{},
			Status:          model.OrderDtoStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.
		On("GetOrder", expectedOrderUUID).
		Return(expectedOrder).
		Once()

	s.paymentClient.
		On("PayOrder", s.ctx, expectedOrder.OrderUUID, expectedOrder.UserUUID, string(expectedPaymentMethod)).
		Return(expectedTransactionUUID, nil).
		Once()

	s.orderRepository.On("UpdateOrder", mock.Anything)

	transactionUUID, err := s.service.PayOrder(s.ctx, expectedPaymentMethod, expectedOrderUUID)
	s.Require().NoError(err)
	s.Require().Len(transactionUUID, 36)
	s.Require().Equal(expectedTransactionUUID, transactionUUID)
}

func (s *ServiceSuite) TestPayOrder_RepoError() {
	expectedPaymentMethod := model.PaymentMethodCARD

	s.orderRepository.On("GetOrder", "not-exist").Return(nil).Once()

	transactionUUID, err := s.service.PayOrder(s.ctx, expectedPaymentMethod, "not-exist")
	s.Require().Error(err)
	var notFoundError *model.NotFoundError
	s.Require().ErrorAs(err, &notFoundError)
	s.Require().Contains(err.Error(), "not-exist")
	s.Require().Empty(transactionUUID)
}

func (s *ServiceSuite) TestPayOrder_PaymentServiceError() {
	var (
		expectedError     = gofakeit.Error()
		expectedOrderUUID = gofakeit.UUID()
		userUUID          = gofakeit.UUID()
		total             = 1800.0
		partUuids         = []string{
			"11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222",
		}

		expectedPaymentMethod = model.PaymentMethodCARD

		expectedOrder = &model.Order{
			OrderUUID:       expectedOrderUUID,
			UserUUID:        userUUID,
			PartUuids:       partUuids,
			TotalPrice:      total,
			TransactionUUID: model.OptNilString{},
			PaymentMethod:   &model.NilOrderDtoPaymentMethod{},
			Status:          model.OrderDtoStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.
		On("GetOrder", expectedOrderUUID).
		Return(expectedOrder).
		Once()

	s.paymentClient.
		On("PayOrder", s.ctx, expectedOrder.OrderUUID, expectedOrder.UserUUID, string(expectedPaymentMethod)).
		Return("", expectedError).
		Once()

	transactionUUID, err := s.service.PayOrder(s.ctx, expectedPaymentMethod, expectedOrderUUID)
	s.Require().Error(err)
	var internalError *model.InternalServerError
	s.Require().ErrorAs(err, &internalError)
	s.Require().Contains(err.Error(), expectedError.Error())
	s.Require().Empty(transactionUUID)
}
