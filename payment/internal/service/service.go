package payment

import (
	"context"
)

type PaymentServise interface {
	PayOrder(_ context.Context, orderID string, userID string, paymentMethod int32) (string, error)
}
