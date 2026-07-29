package inventory

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	apiInventoryV1 "github.com/LushnikovSR/spaceship_factory/inventory/internal/api/inventory/v1"
	config "github.com/LushnikovSR/spaceship_factory/inventory/internal/config"
	repository "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository"
	partRepository "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/part"
	service "github.com/LushnikovSR/spaceship_factory/inventory/internal/service"
	partService "github.com/LushnikovSR/spaceship_factory/inventory/internal/service/part"
	closer "github.com/LushnikovSR/spaceship_factory/platform/pkg/closer"
	inventoryV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository repository.PartRepository
	mongoDBClient       *mongo.Client
	mongoDBHandler      *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = apiInventoryV1.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = partService.NewService(d.PartRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.inventoryRepository == nil {
		d.inventoryRepository = partRepository.NewRepository(d.MongoDBHandler(ctx))
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBHandler(ctx context.Context) *mongo.Database {
	if d.mongoDBHandler == nil {
		d.mongoDBHandler = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandler
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %s", err.Error()))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}
