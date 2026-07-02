package order

import (
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/order/internal/repository/converter"
)

// CreateOrder создает новую запись о заказе.
// Если заказа уже есть в хранилище, возвращает ошибку.
func (r *repository) CreateOrder(order *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[order.OrderUUID]; ok {
		return &model.ConflictError{
			BaseError: model.BaseError{
				Code:    409,
				Message: "The order already exists",
			},
		}
	}

	repoOrder := converter.OrderModelToRepoModel(order)

	r.data[order.OrderUUID] = *repoOrder
	return nil
}
