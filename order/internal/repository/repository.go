package repository

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, uuid string) (*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
}
