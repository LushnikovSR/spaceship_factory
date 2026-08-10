package app

import (
	"context"

	apiPaymentV1 "github.com/LushnikovSR/spaceship_factory/payment/internal/api/payment/v1"
	service "github.com/LushnikovSR/spaceship_factory/payment/internal/service"
	servicePayment "github.com/LushnikovSR/spaceship_factory/payment/internal/service/payment"
	paymentV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = apiPaymentV1.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = servicePayment.NewService()
	}

	return d.paymentService
}
