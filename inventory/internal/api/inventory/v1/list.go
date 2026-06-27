package inventory

import (
	"context"

	converter "github.com/LushnikovSR/spaceship_factory/inventory/internal/converter"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	parts, err := a.inventoryService.ListParts(ctx, converter.RequestToModelPart(req))
	if err != nil {
		return &inventory_v1.ListPartsResponse{}, nil
	}

	return &inventory_v1.ListPartsResponse{
		Parts: converter.ModelListPartsToProto(parts),
	}, nil
}
