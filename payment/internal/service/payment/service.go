package payment

import (
	def "github.com/LushnikovSR/spaceship_factory/payment/internal/service"
)

var _ def.PaymentServise = (*service)(nil)

type service struct {
}

func NewService() *service {
	return &service{}
}
