//go:build integration

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

var (
	testPart = repoModel.Part{
		Name:          "Testing part",
		Description:   "Testing description",
		Price:         4999.89,
		StockQuantity: 19,
		Category:      repoModel.Category(0),
		Dimensions: &repoModel.Dimensions{
			Length: 10,
			Width:  20,
			Height: 30,
			Weight: 40,
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    "TestName",
			Country: "TestCountry",
			Website: "testwebsite.test",
		},
		Tags:      []string{"test", "main"},
		Metadata:  make(map[string]interface{}),
		CreatedAt: time.Now(),
	}
)

// InsertTestPart - вставляет тестовую запчасть в коллекцию Mongo и возвращает UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partDoc := bson.M{
		"name":           testPart.Name,
		"description":    testPart.Description,
		"price":          testPart.Price,
		"stock_quantity": testPart.StockQuantity,
		"category":       testPart.Category,
		"dimensions": bson.M{
			"length": testPart.Dimensions.Length,
			"width":  testPart.Dimensions.Width,
			"height": testPart.Dimensions.Height,
			"weight": testPart.Dimensions.Weight,
		},
		"manufacturer": bson.M{
			"name":    testPart.Manufacturer.Name,
			"country": testPart.Manufacturer.Country,
			"website": testPart.Manufacturer.Website,
		},
		"tags":       testPart.Tags,
		"created_at": testPart.CreatedAt,
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

// InsertTestPartWithData — вставляет тестовую запчасть с заданными данными в коллекцию Mongo и возвращает uuid
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part repoModel.Part) (string, error) {
	now := time.Now()
	var dimention map[string]interface{}

	if part.Dimensions != nil {
		dimention = bson.M{
			"length": "",
			"width":  "",
			"height": part.Dimensions.Height,
			"weight": part.Dimensions.Weight,
		}
	}

	var manufacturer map[string]interface{}

	if part.Manufacturer != nil {
		manufacturer = bson.M{
			"name":    part.Manufacturer.Name,
			"country": part.Manufacturer.Country,
			"website": part.Manufacturer.Website,
		}
	}

	partDoc := bson.M{
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       part.Category,
		"dimensions":     dimention,
		"manufacturer":   manufacturer,
		"tags":           part.Tags,
		"created_at":     primitive.NewDateTimeFromTime(now),
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
func (env *TestEnvironment) GetTestPart() *repoModel.Part {
	return &repoModel.Part{
		Name:          testPart.Name,
		Description:   testPart.Description,
		Price:         testPart.Price,
		StockQuantity: testPart.StockQuantity,
		Category:      testPart.Category,
		Dimensions:    testPart.Dimensions,
		Manufacturer:  testPart.Manufacturer,
		Tags:          testPart.Tags,
		Metadata:      testPart.Metadata,
		CreatedAt:     testPart.CreatedAt,
	}
}

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
