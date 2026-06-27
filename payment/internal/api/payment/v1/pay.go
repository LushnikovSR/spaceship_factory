package payment

import (
	"context"

	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	transactionUuid, err := a.paymentService.PayOrder(ctx, req.OrderUuid, req.UserUuid, int32(req.PaymentMethod))
	if err != nil {
		return nil, err
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}
