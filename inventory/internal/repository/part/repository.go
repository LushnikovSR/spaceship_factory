package inventory

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository"
	"github.com/joho/godotenv"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	data *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		data: db.Collection("parts"),
	}
}

func ConnectMongo(ctx context.Context) (*mongo.Client, error) {
	err := godotenv.Load("deploy/compose/inventory/.env") // вызов из корня проекта
	if err != nil {
		slog.Warn("failed to load .env file", "error", err)
	}

	// Получаем строку подключения из переменной окружения
	dbURI := os.Getenv("MONGO_URI")
	if dbURI == "" {
		slog.Warn("mongo uri not specified")
	}

	// Создаем клиент MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI).
		SetConnectTimeout(5*time.Second).
		SetServerSelectionTimeout(5*time.Second))
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return &mongo.Client{}, fmt.Errorf("failed to connect to database: %v\n", err)
	}

	return client, nil
}
