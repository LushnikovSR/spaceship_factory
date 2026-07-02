package order

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error)
}
