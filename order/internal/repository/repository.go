package repository

import model "github.com/LushnikovSR/spaceship_factory/order/internal/model"

type OrderRepository interface {
	CreateOrder(order *model.Order) error
	GetOrder(uuid string) *model.Order
	UpdateOrder(order *model.Order)
}
