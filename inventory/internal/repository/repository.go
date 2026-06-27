package inventory

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
)

type PartRepository interface {
	GetPart(_ context.Context, uuid string) (*model.Part, error)
	ListParts(_ context.Context, filter *model.PartsFilter) ([]*model.Part, error)
}
