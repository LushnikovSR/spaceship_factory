package inventory

import (
	"context"
	"fmt"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/converter"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartsFilter) (parts []*model.Part, err error) {
	mongoFilter, err := buildMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	repoParts, err := r.find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}

	return converter.ToModelParts(repoParts), nil
}

func buildMongoFilter(filter *model.PartsFilter) (bson.M, error) {
	if filter == nil {
		return bson.M{}, nil
	}

	mongoFilter := bson.M{}

	// UUID
	if len(filter.Uuids) > 0 {
		repoIDs, err := converter.UuidsToRepo(filter.Uuids)
		if err != nil {
			return nil, fmt.Errorf("invalid uuids: %w", err)
		}

		mongoFilter["_id"] = bson.M{
			"$in": repoIDs,
		}
	}

	// Имена
	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{
			"$in": filter.Names,
		}
	}

	// Категории
	if len(filter.Categories) > 0 {
		categories := make([]string, len(filter.Categories))
		for i, c := range filter.Categories {
			categories[i] = string(c)
		}

		mongoFilter["category"] = bson.M{
			"$in": categories,
		}
	}

	// Страны производителя
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{
			"$in": filter.ManufacturerCountries,
		}
	}

	// Теги
	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{
			"$in": filter.Tags,
		}
	}

	return mongoFilter, nil
}

func (r *repository) find(ctx context.Context, filter interface{}) ([]repoModel.Part, error) {
	cursor, err := r.data.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var parts []repoModel.Part
	if err := cursor.All(ctx, &parts); err != nil {
		return nil, err
	}

	return parts, nil
}
