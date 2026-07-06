package payment

import (
	"context"
)

type PaymentService interface {
	PayOrder(_ context.Context, orderID, userID string, paymentMethod int32) (string, error)
}
