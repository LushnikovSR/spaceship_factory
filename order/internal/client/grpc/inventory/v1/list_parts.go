package order

import (
	"context"

	clientConverter "github.com/LushnikovSR/spaceship_factory/order/internal/client/converter"
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	resp, err := c.generatedClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	return clientConverter.PartsListToModel(resp.Parts), nil

}
