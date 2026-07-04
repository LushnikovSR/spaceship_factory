package payment

import (
	service "github.com/LushnikovSR/spaceship_factory/payment/internal/service"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
