package order

import (
	"github.com/brianvoe/gofakeit/v7"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

func (s *ServiceSuite) TestGetOrder_Success() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		partUuids = []string{
			"11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222",
		}

		total = 1800.0

		expectOrder = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       partUuids,
			TotalPrice:      total,
			TransactionUUID: model.OptNilString{},
			PaymentMethod:   &model.NilOrderDtoPaymentMethod{},
			Status:          model.OrderDtoStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.
		On("GetOrder", orderUUID).
		Return(&expectOrder).
		Once()

	order, err := s.service.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	// Сравниваем поля по отдельности
	s.Require().Equal(expectOrder.OrderUUID, order.OrderUUID)
	s.Require().Equal(expectOrder.UserUUID, order.UserUUID)
	s.Require().Equal(expectOrder.PartUuids, order.PartUuids)
	s.Require().Equal(expectOrder.TotalPrice, order.TotalPrice)
	s.Require().Equal(expectOrder.Status, order.Status)

	// Для TransactionUUID (значение) сравниваем прямо
	s.Require().Equal(expectOrder.TransactionUUID, order.TransactionUUID)

	// Для указателя PaymentMethod сравниваем содержимое, если не nil
	if expectOrder.PaymentMethod != nil {
		s.Require().NotNil(order.PaymentMethod)
		s.Require().Equal(*expectOrder.PaymentMethod, *order.PaymentMethod)
	} else {
		s.Require().Nil(order.PaymentMethod)
	}
}

func (s *ServiceSuite) TestGetOrder_NotFoundError() {
	s.orderRepository.
		On("GetOrder", "not-exist").
		Return(nil).
		Once()

	order, err := s.service.GetOrder(s.ctx, "not-exist")

	s.Require().Error(err)
	var notFoundError *model.NotFoundError
	s.Require().ErrorAs(err, &notFoundError)
	s.Require().Empty(order)
}
