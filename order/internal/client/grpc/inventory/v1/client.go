package order

import (
	def "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient inventory_v1.InventoryServiceClient
}

func NewClient(generatedClient inventory_v1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
