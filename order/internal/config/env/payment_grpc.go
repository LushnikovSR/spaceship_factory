package order

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type paymentGRPCClientEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type paymentGRPCClientConfig struct {
	raw paymentGRPCClientEnvConfig
}

func NewPaymentGRPCClientConfig() (*paymentGRPCClientConfig, error) {
	var raw paymentGRPCClientEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}
	return &paymentGRPCClientConfig{raw: raw}, nil
}

func (cfg *paymentGRPCClientConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
