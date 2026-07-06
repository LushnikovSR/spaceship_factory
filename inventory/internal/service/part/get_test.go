package inventory

import (
	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestGetPart_Success() {
	var (
		orderUUID = gofakeit.UUID()

		expectedPart = &model.Part{
			UUID: orderUUID,
		}
	)

	s.partRepository.
		On("GetPart", s.ctx, orderUUID).
		Return(expectedPart, nil).
		Once()

	part, err := s.service.GetPart(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(expectedPart, part)
	s.Require().Len(part.UUID, 36)
	s.Require().Equal(part.UUID, orderUUID)
}

func (s *ServiceSuite) TestGetPart_RepoError() {
	var (
		expectedError = gofakeit.Error()

		nonExistentOrderUUID = "non-existent"
	)

	s.partRepository.
		On("GetPart", s.ctx, nonExistentOrderUUID).
		Return(&model.Part{}, expectedError).
		Once()

	part, err := s.service.GetPart(s.ctx, nonExistentOrderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, expectedError)
	s.Require().Empty(part)
}
