//go:build integration

package integration

import (
	"context"

	"github.com/LushnikovSR/spaceship_factory/platform/pkg/logger"
	"go.uber.org/zap"
)

//teardownTestEnvironment - освобождает все ресурсы тестового окружения
func teardownTestEnvironment(ctx context.Context, env *TestEnvironment) {
	log := logger.Logger()
	log.Info(ctx, "🧹 cleaning up the test environment")

	cleanupTestEnvironment(ctx, env)

	log.Info(ctx, "✅ test environment has been successfully cleaned up")
}

// cleanupTestEnvironment - вспомогательная функция для осовбождения ресурсов
func cleanupTestEnvironment(ctx context.Context, env *TestEnvironment) {
	if env.App != nil {
		err := env.App.Terminate(ctx)
		if err != nil {
			logger.Error(ctx, "failed to terminate App container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 App container successfuly stopped")
		}
	}

	if env.Mongo != nil {
		err := env.Mongo.Terminate(ctx)
		if err != nil {
			logger.Error(ctx, "failed to terminate MongoDb container", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 MongoDB container successfuly stopped")
		}
	}

	if env.Network != nil {
		err := env.Network.Remove(ctx)
		if err != nil {
			logger.Error(ctx, "failed to remove container network", zap.Error(err))
		} else {
			logger.Info(ctx, "🛑 Network successfully removed")
		}
	}
}
