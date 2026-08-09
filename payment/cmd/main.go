package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/LushnikovSR/spaceship_factory/payment/internal/app"
	config "github.com/LushnikovSR/spaceship_factory/payment/internal/config"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/closer"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/logger"
)

const (
	shutdownTimeout = 5 * time.Second
	configPath      = "./deploy/compose/payment/.env"
)

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefullShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}

	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
	}
}

func gracefullShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := closer.CloseAll(ctx)
	if err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
