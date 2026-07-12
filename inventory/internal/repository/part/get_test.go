package inventory

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
)

func (s *RepositorySuite) TestGetPart_Success() {
	var (
		existedPartUUID = "11111111-1111-1111-1111-111111111111"
		expectedPart    = &model.Part{
			UUID:          "11111111-1111-1111-1111-111111111111",
			Name:          "Сопло маршевое",
			Price:         1500.0,
			StockQuantity: 5,
			Category:      model.CATEGORY_ENGINE,
			Manufacturer: &model.Manufacturer{
				Name:    "Biscuit",
				Country: "Germany",
				Website: "financialharness.info",
			},
			Tags: []string{"engine", "main"},
		}
	)
	s.repository.Init(context.Background())

	part, err := s.repository.GetPart(s.ctx, existedPartUUID)
	s.Require().NoError(err)
	s.Require().Equal(expectedPart, part)
}

func (s *RepositorySuite) TestGetPart_NotFoundError() {
	nonExistentPartUUID := "non-existent"

	s.repository.Init(context.Background())

	part, err := s.repository.GetPart(s.ctx, nonExistentPartUUID)
	s.Require().Error(err)
	s.Require().ErrorAs(err, &model.ErrPartNotFound)
	s.Require().Empty(part)
}
