package order

import (
	"context"

	"github.com/google/uuid"
)

// NewOrder creates order.
func (a *api) NewOrder(ctx context.Context, partUuids []string) (orderUUID string, totalPrice float64, err error) {
	userUUID := uuid.NewString()
	if len(partUuids) == 0 {
		partUuids = []string{
			"11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222",
		}
	}

	orderUUID, totalPrice, err = a.orderService.CreateOrder(ctx, userUUID, partUuids)
	if err != nil {
		return "", 0, err
	}
	return orderUUID, totalPrice, nil
}
