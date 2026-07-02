package order

import (
	def "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient payment_v1.PaymentServiceClient
}

func NewClient(generatedClient payment_v1.PaymentServiceClient) *client {
	return &client{generatedClient: generatedClient}
}
