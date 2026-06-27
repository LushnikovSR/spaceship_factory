package inventory

import (
	repository "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository"
	def "github.com/LushnikovSR/spaceship_factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
