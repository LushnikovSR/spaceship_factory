package integration

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
)

// InsertTestPart - вставляет тестовую запчасть в коллекцию Mongo и возвращает UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	now := time.Now()

	partDoc := bson.M{
		"name":           "Testing part",
		"description":    "Testing description",
		"price":          4999.89,
		"stock_quantity": 19,
		"category":       repoModel.CATEGORY_UNSPECIFIED,
		"dimensions": bson.M{
			"length": 10,
			"width":  20,
			"height": 30,
			"weight": 40,
		},
		"manufacturer": bson.M{
			"name":    "TestName",
			"country": "TestCounty",
			"website": "testwebsite.test",
		},
		"tags":       []string{"engine", "main"},
		"created_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	inserPart, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	var idString string

	if oid, ok := inserPart.InsertedID.(primitive.ObjectID); ok {
		idString = oid.Hex()
	} else {
		return "", fmt.Errorf("failed to convert ID to ObjectID: %w", err)
	}

	return idString, nil
}

// InsertTestPartWithData — вставляет тестовое запчасть с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part repoModel.Part) (string, error) {
	now := time.Now()

	partDoc := bson.M{
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       part.Category,
		"dimensions": bson.M{
			"length": part.Dimensions.Length,
			"width":  part.Dimensions.Width,
			"height": part.Dimensions.Height,
			"weight": part.Dimensions.Weight,
		},
		"manufacturer": bson.M{
			"name":    part.Manufacturer.Name,
			"country": part.Manufacturer.Country,
			"website": part.Manufacturer.Website,
		},
		"tags":       part.Tags,
		"created_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	inserPart, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	var idString string

	if oid, ok := inserPart.InsertedID.(primitive.ObjectID); ok {
		idString = oid.Hex()
	} else {
		return "", fmt.Errorf("failed to convert ID to ObjectID: %w", err)
	}

	return idString, nil
}

// GetTestPartInfo — возвращает тестовую информацию о запчасте

// GetUpdatedParInfo — возвращает обновленную информацию о запчасте

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil

}
