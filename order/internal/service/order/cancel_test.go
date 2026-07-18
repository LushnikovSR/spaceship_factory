package order

import (
	"github.com/brianvoe/gofakeit/v7"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
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
		On("GetOrder", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	s.orderRepository.On("UpdateOrder", s.ctx, order).
		Return(nil).
		Once()

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrder_RepoError_NotFound() {
	s.orderRepository.
		On("GetOrder", s.ctx, "non-existent").
		Return(&model.Order{}, nil).
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
		On("GetOrder", s.ctx, orderUUID).
		Return(order, nil).
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
		On("GetOrder", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	err := s.service.CancelOrder(s.ctx, orderUUID)
	s.Require().Error(err)
	var badRequestError *model.BadRequestError
	s.Require().ErrorAs(err, &badRequestError)
	s.Require().Contains(err.Error(), orderUUID)
}
