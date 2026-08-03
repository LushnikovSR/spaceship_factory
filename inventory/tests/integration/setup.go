package integration

import (
	"context"
	"os"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"

	"github.com/LushnikovSR/spaceship_factory/platform/pkg/logger"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/testcontainers"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/testcontainers/app"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/testcontainers/mongo"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/testcontainers/network"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/testcontainers/path"
)

var (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Переменные окружения приложения
	grpcPortKey = "GRPC_PORT"

	// Значение переменных окружения
	loggerLevelValue = "debug"
	startTimeout     = 3 * time.Minute
)

// TestEnvironment - структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// setupTestEnvironment - подготавливает тестовое окружение: сеть, контейнеры и возвращает структуру с ресурсами
func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Preparing the test environment...")

	// Создаётся общая docker-сеть
	generatedNetwork, err := network.NewNetwork(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "failed to creat common docker network", zap.Error(err))
	}
	logger.Info(ctx, "✅ Common docker network was successfully created")

	// Получаем переменные окружения для MongoDB с проверкой на наличие
	mongoUsername := getEnvWithLogging(ctx, testcontainers.MongoUsernameKey)
	mongoPassword := getEnvWithLogging(ctx, testcontainers.MongoPasswordKey)
	mongoImageName := getEnvWithLogging(ctx, testcontainers.MongoImageNameKey)
	mongoDatabase := getEnvWithLogging(ctx, testcontainers.MongoDatabaseKey)

	// Получаем порт gRPC для waitStrategy
	grpcPort := getEnvWithLogging(ctx, grpcPortKey)

	// Запускаем контейнер с MongoDB
	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoConteinerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ MongoDB container has started successfully")

	// Запускаем контейнер с приложением
	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		testcontainers.MongoHostKey: generatedMongo.Config().ContainerName,
	}

	// Создаём настраиваемую стратегию с увеличенным таймаутом
	withStratagy := wait.ForListeningPort(string(nat.Port(grpcPort + "/tcp"))).
		WithStartupTimeout(startTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithPort(grpcPort),
		app.WithEnv(appEnv),
		app.WithNetworks(generatedNetwork.Name()),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(withStratagy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "application container could not be started", zap.Error(err))
	}

	logger.Info(ctx, "✅ application container has started successfully")

	logger.Info(ctx, "🎉 The test environment is ready")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

// getEnvWithLogging - возвращает значение переменной окружения с логгированием
func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Environment variable value was not set", zap.String("key", key))
	}

	return value
}
