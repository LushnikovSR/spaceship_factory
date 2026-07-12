package inventory

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/converter"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
)

// GetByID gets a note by ID
func (r *repository) GetPart(ctx context.Context, id string) (*model.Part, error) {
	var part repoModel.Part

	objID, err := primitive.ObjectIDFromHex(id) // ObjectIDFromHex takes an argument of 24 characters in length
	if err != nil {
		return &model.Part{}, fmt.Errorf("Неверный формат строки: %v", err)
	}

	err = r.data.FindOne(ctx, bson.M{"_id": objID}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &model.Part{}, model.ErrPartNotFound
		}

		return &model.Part{}, err
	}

	return converter.RepoModelToModel(part), nil
}
