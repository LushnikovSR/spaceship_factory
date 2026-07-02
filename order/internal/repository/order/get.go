package order

import (
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/order/internal/repository/converter"
)

// GetOrder возвращает информацию о заказе по uuid.
// Если заказ не найден, возвращает nil.
func (r *repository) GetOrder(uuid string) *model.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[uuid]
	if !ok {
		return nil
	}

	modelOrder := converter.OrderRepoModelToModel(&order)

	return modelOrder
}
