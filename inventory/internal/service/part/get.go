package inventory

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.partRepository.GetPart(ctx, uuid)
	if err != nil {
		return &model.Part{}, err
	}
	return part, nil
}
