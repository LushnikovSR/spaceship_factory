package inventory

import (
	"context"

	converter "github.com/LushnikovSR/spaceship_factory/inventory/internal/converter"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	part, err := a.inventoryService.GetPart(ctx, req.Uuid)
	if err != nil {
		return &inventory_v1.GetPartResponse{}, err
	}

	protoPart := converter.ModelToProto(part)

	return &inventory_v1.GetPartResponse{
		Part: protoPart,
	}, nil
}
