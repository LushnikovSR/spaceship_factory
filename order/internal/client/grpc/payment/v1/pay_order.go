package order

import (
	"context"

	clientConverter "github.com/LushnikovSR/spaceship_factory/order/internal/client/converter"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID, paymentMethod string) (string, error) {
	method, err := clientConverter.PaymentMethodToProto(paymentMethod)
	if err != nil {
		return "", err
	}

	request := &payment_v1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: *method,
	}
	resp, err := c.generatedClient.PayOrder(ctx, request)
	if err != nil {
		return "", nil
	}
	return resp.TransactionUuid, nil
}
