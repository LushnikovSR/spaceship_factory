package inventory

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (*model.Part, error) {
	part := r.Read(uuid)
	if part == nil {
		return &model.Part{}, fmt.Errorf("part uuid %s: %w", uuid, model.ErrPartNotFound)
	}

	return part, nil
}

func (r *repository) Read(uuid string) *model.Part {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[uuid]
	if !ok {
		return nil
	}
	return converter.RepoModelToModel(lo.ToPtr(part))
}
