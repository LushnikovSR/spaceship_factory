package inventory

import (
	inventory "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelToProto(inmodel *model.Part) *inventory_v1.Part {
	if inmodel == nil {
		return nil
	}
	metadata := make(map[string]*inventory_v1.Value, len(inmodel.Metadata))
	for k, v := range inmodel.Metadata {
		metadata[k] = RepoValueToModelValue(v)
	}

	dimensions := &inventory_v1.Dimensions{
		Length: inmodel.Dimensions.Length,
		Width:  inmodel.Dimensions.Width,
		Height: inmodel.Dimensions.Height,
		Weight: inmodel.Dimensions.Weight,
	}

	manufacturer := &inventory_v1.Manufacturer{
		Name:    inmodel.Manufacturer.Name,
		Country: inmodel.Manufacturer.Country,
		Website: inmodel.Manufacturer.Website,
	}

	createdAt := timestamppb.New(*inmodel.CreatedAt)
	updatedAt := timestamppb.New(*inmodel.UpdatedAt)

	return &inventory_v1.Part{
		Uuid:          inmodel.UUID,
		Name:          inmodel.Name,
		Description:   inmodel.Description,
		Price:         inmodel.Price,
		StockQuantity: inmodel.StockQuantity,
		Category:      inventory_v1.Category(inmodel.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          inmodel.Tags,
		Metadata:      metadata,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func RepoValueToModelValue(v *model.Value) *inventory_v1.Value {
	if v == nil {
		return nil
	}
	mv := &inventory_v1.Value{}
	switch data := v.DataType.(type) {
	case *model.Value_StringValue:
		mv.DataType = &inventory_v1.Value_StringValue{StringValue: data.StringValue}
	case *model.Value_Int64Value:
		mv.DataType = &inventory_v1.Value_Int64Value{Int64Value: data.Int64Value}
	case *model.Value_DoubleValue:
		mv.DataType = &inventory_v1.Value_DoubleValue{DoubleValue: data.DoubleValue}
	case *model.Value_BoolValue:
		mv.DataType = &inventory_v1.Value_BoolValue{BoolValue: data.BoolValue}
	}
	return mv
}

func RequestToModelPart(req *inventory_v1.ListPartsRequest) *model.PartsFilter {
	categories := make([]model.Category, 0, len(req.Filter.Categories))

	for _, category := range req.Filter.Categories {
		categories = append(categories, model.Category(category))
	}

	return &inventory.PartsFilter{
		Uuids:                 req.Filter.Uuids,
		Names:                 req.Filter.Names,
		Categories:            categories,
		ManufacturerCountries: req.Filter.ManufacturerCountries,
		Tags:                  req.Filter.Tags,
	}
}

func ModelListPartsToProto(inParts []*inventory.Part) []*inventory_v1.Part {
	parts := make([]*inventory_v1.Part, 0, len(inParts))

	for _, part := range inParts {
		protoPart := ModelToProto(part)
		parts = append(parts, protoPart)
	}
	return parts
}
