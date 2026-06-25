package payment

import (
	def "github.com/LushnikovSR/spaceship_factory/payment/internal/service/service.go"
)

var _ def.PaymentService = (*service)(nil)

type service struct {
}

func NewService() *service {
	return &service{}
}
