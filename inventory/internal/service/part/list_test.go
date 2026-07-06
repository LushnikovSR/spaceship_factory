package inventory

import (
	"github.com/brianvoe/gofakeit/v7"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListParts_Success() {
	var (
		firstPartUUID  = gofakeit.UUID()
		secondPartUUID = gofakeit.UUID()

		uuids = []string{firstPartUUID, secondPartUUID}

		filter = &model.PartsFilter{
			Uuids: uuids,
		}

		expectedParts = []*model.Part{{UUID: firstPartUUID}, {UUID: secondPartUUID}}
	)

	s.partRepository.
		On("ListParts", s.ctx, filter).
		Return(expectedParts, nil).
		Once()

	parts, err := s.service.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedParts[0].UUID, parts[0].UUID)
	s.Require().Equal(expectedParts[1].UUID, parts[1].UUID)
}

func (s *ServiceSuite) TestListParts_RepoError_NonExistentPartID() {
	var (
		expectedError = gofakeit.Error()

		firstPartUUID = gofakeit.UUID()

		uuids = []string{firstPartUUID, "non-existent"}

		filter = &model.PartsFilter{
			Uuids: uuids,
		}
	)

	s.partRepository.
		On("ListParts", s.ctx, filter).
		Return([]*model.Part{}, expectedError).
		Once()

	parts, err := s.service.ListParts(s.ctx, filter)
	s.Require().Error(err)
	s.Require().Empty(parts)
}

func (s *ServiceSuite) TestListParts_NilFilter() {
	var (
		filter           = (*model.PartsFilter)(nil)
		expectedAllParts = []*model.Part{}
	)
	s.partRepository.
		On("ListParts", s.ctx, filter).
		Return(expectedAllParts, nil).
		Once()

	parts, err := s.service.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedAllParts, parts)
}
