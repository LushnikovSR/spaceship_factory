package order

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
)

type OrderService interface {
	CancelOrder(ctx context.Context, uuid string) error
	CreateOrder(ctx context.Context, userUUID string, partUuids []string) (orderUUID string, totalPrice float64, err error)
	GetOrder(ctx context.Context, orderUUID string) (model.Order, error)
	PayOrder(ctx context.Context, paymentMethod model.PaymentMethod, orderUUID string) (transactionUUID string, err error)
}
