package payment

import ()

type PaymentServise interface {
	PayOrder(_ context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error)
}
