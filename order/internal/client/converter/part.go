package order

import (
	"fmt"
	"time"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

func PartsFilterToProto(filter model.PartsFilter) *inventory_v1.PartsFilter {
	categories := make([]inventory_v1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventory_v1.Category(category))
	}

	return &inventory_v1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func PartsListToModel(parts []*inventory_v1.Part) []model.Part {
	modelParts := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		modelPart := PartToModel(part)
		modelParts = append(modelParts, model.Part(modelPart))
	}

	return modelParts
}

func PartToModel(part *inventory_v1.Part) model.Part {
	if part == nil {
		return model.Part{}
	}

	metadata := make(map[string]*model.Value, len(part.Metadata))
	if len(part.Metadata) != 0 {
		for k, v := range part.Metadata {
			metadata[k] = ValueToModel(v)
		}
	}

	// Если Dimensions == nil, оставляем поле nil
	var dimensions *model.Dimensions
	if part.Dimensions != nil {
		dimensions = &model.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		}
	}

	// Если Manufacturer == nil, оставляем поле nil
	var manufacturer *model.Manufacturer
	if part.Manufacturer != nil {
		manufacturer = &model.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		}
	}

	// Если createdAt или updatedAt == nil, оставляем nil
	var createdAt, updatedAt *time.Time
	if part.CreatedAt != nil {
		temp := part.CreatedAt.AsTime()
		createdAt = &temp
	}

	if part.UpdatedAt != nil {
		temp := part.UpdatedAt.AsTime()
		updatedAt = &temp
	}

	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          part.Tags,
		Metadata:      metadata,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func ValueToModel(value *inventory_v1.Value) *model.Value {
	if value == nil {
		return nil
	}

	mv := &model.Value{}
	switch data := value.DataType.(type) {
	case *inventory_v1.Value_StringValue:
		mv = &model.Value{DataType: &model.Value_StringValue{StringValue: data.StringValue}}
	case *inventory_v1.Value_Int64Value:
		mv = &model.Value{DataType: &model.Value_Int64Value{Int64Value: data.Int64Value}}
	case *inventory_v1.Value_DoubleValue:
		mv = &model.Value{DataType: &model.Value_DoubleValue{DoubleValue: data.DoubleValue}}
	case *inventory_v1.Value_BoolValue:
		mv = &model.Value{DataType: &model.Value_BoolValue{BoolValue: data.BoolValue}}
	case nil:
		mv = &model.Value{DataType: nil}
	}
	return mv
}

func PaymentMethodToProto(method string) (*payment_v1.PaymentMethod, error) {
	value, ok := PaymentMethod_value[method]
	if !ok {
		return nil, &model.BadRequestError{
			BaseError: model.BaseError{
				Code:    400,
				Message: fmt.Sprintf("Conversion to proto failed, Invalid payment method: %s", method),
			},
		}
	}
	return (*payment_v1.PaymentMethod)(&value), nil
}

/*
PaymentMethod contains a list of possible payment methods

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)
*/
// Enum value maps for PaymentMethod.
var (
	PaymentMethod_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "CARD",
		2: "SBP",
		3: "CARD",
		4: "MONEY",
	}
	PaymentMethod_value = map[string]int32{
		"UNSPECIFIED":    0,
		"CARD":           1,
		"SBP":            2,
		"CREDIT_CARD":    3,
		"INVESTOR_MONEY": 4,
	}
)
