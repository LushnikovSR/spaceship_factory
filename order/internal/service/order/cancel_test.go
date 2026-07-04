package order

import (
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestCancelOrder_Success() {
	var (
		orderUUID = gofakeit.UUID()

		order = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderDtoStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.
		On("GetOrder", orderUUID).
		Return(order).
		Once()

	s.orderRepository.On("UpdateOrder", order)

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrder_RepoError_NotFound() {
	s.orderRepository.
		On("GetOrder", "non-existent").
		Return(nil).
		Once()

	err := s.service.CancelOrder(s.ctx, "non-existent")
	s.Require().Error(err)
	var notFoundError *model.NotFoundError
	s.Require().ErrorAs(err, &notFoundError)
	s.Require().Contains(err.Error(), "non-existent")
}

func (s *ServiceSuite) TestCancelOrder_RepoError_Conflict() {
	var (
		orderUUID = gofakeit.UUID()

		order = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderDtoStatusPAID,
		}
	)

	s.orderRepository.
		On("GetOrder", orderUUID).
		Return(order).
		Once()

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Require().Error(err)
	var conflictError *model.ConflictError
	s.Require().ErrorAs(err, &conflictError)
	s.Require().Contains(err.Error(), orderUUID)
}

func (s *ServiceSuite) TestCancelOrder_RepoError_BadRequest() {
	var (
		orderUUID = gofakeit.UUID()

		order = &model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderDtoStatusCANCELLED,
		}
	)

	s.orderRepository.
		On("GetOrder", orderUUID).
		Return(order).
		Once()

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Require().Error(err)
	var badRequestError *model.BadRequestError
	s.Require().ErrorAs(err, &badRequestError)
	s.Require().Contains(err.Error(), orderUUID)
}
