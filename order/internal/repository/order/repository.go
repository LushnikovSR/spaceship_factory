package order

import (
	"sync"

	def "github.com/LushnikovSR/spaceship_factory/order/internal/repository"
	repoModel "github.com/LushnikovSR/spaceship_factory/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Order),
	}
}
