package order

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryGRPCClintEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type inventoryGRPCClientConfig struct {
	raw inventoryGRPCClintEnvConfig
}

func NewInventoryGRPCClientConfig() (*inventoryGRPCClientConfig, error) {
	var raw inventoryGRPCClintEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &inventoryGRPCClientConfig{raw: raw}, nil
}

func (cfg *inventoryGRPCClientConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
