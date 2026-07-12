package inventory

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *repository) Init(ctx context.Context) {
	// Наполняем тестовыми данными
	id1, err := r.Create(ctx, repoModel.Part{
		Name:          "Сопло маршевое",
		Price:         1500.0,
		StockQuantity: 5,
		Category:      repoModel.CATEGORY_ENGINE,
		Manufacturer: &repoModel.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		},
		Tags: []string{"engine", "main"},
	})
	if err != nil {
		slog.Warn(fmt.Errorf("added part is failed error: %w", err).Error())
	} else {
		slog.Info(fmt.Sprintf("part with id %s correctly added", id1))
	}

	id2, err := r.Create(ctx, repoModel.Part{
		Name:          "Иллюминатор стандартный",
		Price:         300.0,
		StockQuantity: 12,
		Category:      repoModel.CATEGORY_PORTHOLE,
		Manufacturer: &repoModel.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		},
		Tags: []string{"porthole", "window"},
	})

	if err != nil {
		slog.Warn(fmt.Errorf("added part is failed error: %w", err).Error())
	} else {
		slog.Info(fmt.Sprintf("part with id %s correctly added", id2))
	}

	id3, err := r.Create(ctx, repoModel.Part{
		Name:          "Иллюминатор квадратный",
		Price:         600.0,
		StockQuantity: 2,
		Category:      repoModel.CATEGORY_PORTHOLE,
		Manufacturer:  nil,
		Tags:          nil,
	})

	if err != nil {
		slog.Warn(fmt.Errorf("added part is failed error: %w", err).Error())
	} else {
		slog.Info(fmt.Sprintf("part with id %s correctly added", id3))
	}

}

// Create создает новую запись в mongoDB.Collection
func (r *repository) Create(ctx context.Context, part repoModel.Part) (string, error) {
	if part.CreatedAt.IsZero() {
		part.CreatedAt = time.Now()
	}

	res, err := r.data.InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	val, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to convert _id to string")
	}

	return val.Hex(), nil
}

// Update обновляет запись в mongoDB.Collection
func (r *repository) Update(ctx context.Context, part repoModel.Part) (int, error) {
	// UpdateOne обновляет первый документ, соответствующий фильтру
	// Первый параметр - фильтр для поиска документа (по ID)
	// Второй параметр - операции обновления:
	//   $set - устанавливает новые значения для указанных полей
	updateResult, err := r.data.UpdateOne(
		ctx,
		bson.M{"_id": part.ID},
		bson.M{
			"$set": bson.M{
				"price":          part.Price,
				"stock_quantity": part.StockQuantity,
				"updated_at":     time.Now(),
			},
		},
	)
	if err != nil {
		return 0, err
	}

	return int(updateResult.ModifiedCount), nil
}
