package inventory

import (
	service "github.com/LushnikovSR/spaceship_factory/inventory/internal/service"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewAPI(inventoryService service.InventoryService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
