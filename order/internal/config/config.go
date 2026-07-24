package order

import (
	"os"

	"github.com/joho/godotenv"

	env "github.com/LushnikovSR/spaceship_factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger        LoggerConfig
	Postgres      PostgresConfig
	InventoryGRPC InventoryGRPCConfig
	PaymentGRPC   PaymentGRPCConfig
	OrderHTTP     OrderHTTPConfig
}

func Load(paths ...string) error {
	err := godotenv.Load(paths...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostrgresCOnfig()
	if err != nil {
		return err
	}

	inventoryGPRCCfg, err := env.NewInventoryGRPCClientConfig()
	if err != nil {
		return err
	}

	paymentGPRCCfg, err := env.NewPaymentGRPCClientConfig()
	if err != nil {
		return err
	}

	orderCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:        loggerCfg,
		Postgres:      postgresCfg,
		InventoryGRPC: inventoryGPRCCfg,
		PaymentGRPC:   paymentGPRCCfg,
		OrderHTTP:     orderCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
