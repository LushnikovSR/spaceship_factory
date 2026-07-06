package inventory

import (
	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
)

func RepoModelToModel(inmodel *repoModel.Part) *model.Part {
	if inmodel == nil {
		return nil
	}

	// Обрабатываем Metadata: если nil, то и в результате nil
	var metadata map[string]*model.Value
	if inmodel.Metadata != nil {
		metadata = make(map[string]*model.Value, len(inmodel.Metadata))
		for k, v := range inmodel.Metadata {
			metadata[k] = RepoValueToModelValue(v)
		}
	}

	return &model.Part{
		UUID:          inmodel.UUID,
		Name:          inmodel.Name,
		Description:   inmodel.Description,
		Price:         inmodel.Price,
		StockQuantity: inmodel.StockQuantity,
		Category:      model.Category(inmodel.Category),
		Dimensions:    (*model.Dimensions)(inmodel.Dimensions),
		Manufacturer:  (*model.Manufacturer)(inmodel.Manufacturer),
		Tags:          inmodel.Tags,
		Metadata:      metadata,
		CreatedAt:     inmodel.CreatedAt,
		UpdatedAt:     inmodel.UpdatedAt,
	}
}

func RepoValueToModelValue(v *repoModel.Value) *model.Value {
	if v == nil {
		return nil
	}
	mv := &model.Value{}
	switch data := v.DataType.(type) {
	case *repoModel.Value_StringValue:
		mv.DataType = &model.Value_StringValue{StringValue: data.StringValue}
	case *repoModel.Value_Int64Value:
		mv.DataType = &model.Value_Int64Value{Int64Value: data.Int64Value}
	case *repoModel.Value_DoubleValue:
		mv.DataType = &model.Value_DoubleValue{DoubleValue: data.DoubleValue}
	case *repoModel.Value_BoolValue:
		mv.DataType = &model.Value_BoolValue{BoolValue: data.BoolValue}
	}
	return mv
}
