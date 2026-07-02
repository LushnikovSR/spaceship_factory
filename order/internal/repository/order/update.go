package order

import (
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/order/internal/repository/converter"
)

// UpdateOrder обновляет данные о заказе для указанного заказа.
// Если заказа нет в хранилище, создает новую запись.
func (r *repository) UpdateOrder(order *model.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	repoOrder := converter.OrderModelToRepoModel(order)

	r.data[order.OrderUUID] = *repoOrder
}
