package payment

import (
	model "github.com/LushnikovSR/spaceship_factory/payment/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestPayOrder_Success() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = int32(1)
	)

	transactionUuid, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, paymentMethod)
	s.Require().NoError(err)
	s.Require().NotEmpty(transactionUuid)
	s.Require().Len(transactionUuid, 36)
}

func (s *ServiceSuite) TestPayOrder_MissingOrderIDError() {
	var (
		userUUID      = gofakeit.UUID()
		paymentMethod = int32(1)
	)

	transactionUuid, err := s.service.PayOrder(s.ctx, "", userUUID, paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrMissingArgument)
	s.Require().Empty(transactionUuid)
}

func (s *ServiceSuite) TestPayOrder_MissingUserIDError() {
	var (
		orderUUID     = gofakeit.UUID()
		paymentMethod = int32(1)
	)

	transactionUuid, err := s.service.PayOrder(s.ctx, orderUUID, "", paymentMethod)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrMissingArgument)
	s.Require().Empty(transactionUuid)
}

func (s *ServiceSuite) TestPayOrder_MissingPaymentMethodError() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()
	)

	transactionUuid, err := s.service.PayOrder(s.ctx, orderUUID, userUUID, 0)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrMissingArgument)
	s.Require().Empty(transactionUuid)
}
