package inventory

import (
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
)

func RepoModelToModel(part repoModel.Part) *model.Part {
	// Обрабатываем Metadata: если nil, то и в результате nil
	var metadata map[string]*model.Value
	if part.Metadata != nil {
		metadata = make(map[string]*model.Value, len(part.Metadata))
		for k, v := range part.Metadata {
			metadata[k] = RepoValueToModelValue(v)
		}
	}

	return &model.Part{
		UUID:          part.ID.Hex(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    (*model.Dimensions)(part.Dimensions),
		Manufacturer:  (*model.Manufacturer)(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadata,
		CreatedAt:     &part.CreatedAt,
		UpdatedAt:     &part.UpdatedAt,
	}
}

// RepoValueToModelValue конвертирует значение, прочитанное из MongoDB (map[string]interface{}),
// в типизированное представление model.Value.
func RepoValueToModelValue(v any) *model.Value {
	if v == nil {
		return nil
	}
	mv := &model.Value{}
	switch val := v.(type) {
	case string:
		mv.DataType = &model.Value_StringValue{StringValue: val}
	case int64:
		mv.DataType = &model.Value_Int64Value{Int64Value: val}
	case float64:
		mv.DataType = &model.Value_DoubleValue{DoubleValue: val}
	case bool:
		mv.DataType = &model.Value_BoolValue{BoolValue: val}
	default:
		// неизвестный тип — оставляем nil
		slog.Warn("failed to convert to model.Value")
		return nil
	}
	return mv
}

func ToModelParts(repoParts []repoModel.Part) []*model.Part {
	parts := make([]*model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		part := RepoModelToModel(repoPart)
		parts = append(parts, part)
	}
	return parts
}

func UuidsToRepo(uuids []string) ([]primitive.ObjectID, error) {
	repoUuids := make([]primitive.ObjectID, 0, len(uuids))
	for _, uuid := range uuids {
		repoUUID, err := primitive.ObjectIDFromHex(uuid)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to primitive.ObjectID: %w", err)
		}
		repoUuids = append(repoUuids, repoUUID)
	}
	return repoUuids, nil
}
