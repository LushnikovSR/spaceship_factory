package inventory

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

func TestModelToProto_NilInput(t *testing.T) {
	assert.Nil(t, ModelToProto(nil))
}

func TestModelToProto_FullDataSuccess(t *testing.T) {
	var (
		length = gofakeit.Float64Range(1, 1000)
		width  = gofakeit.Float64Range(1, 1000)
		height = gofakeit.Float64Range(1, 1000)
		weight = gofakeit.Float64Range(1, 1000)

		manufacurerName     = gofakeit.ProductName()
		manufacurerCountry  = gofakeit.Country()
		manufacturerWebsite = gofakeit.DomainName()

		partUUID          = gofakeit.UUID()
		partName          = gofakeit.Name()
		partDescription   = gofakeit.HackerPhrase()
		partPrice         = gofakeit.Price(100.0, 2000.0)
		partStockQuantity = gofakeit.RandomInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
		partCategory      = inventory_v1.Category_CATEGORY_ENGINE
		partDimensions    = &inventory_v1.Dimensions{Length: length, Width: width, Height: height, Weight: weight}
		partManufacturer  = &inventory_v1.Manufacturer{Name: manufacurerName, Country: manufacurerCountry, Website: manufacturerWebsite}
		partTags          = []string{}
		metadata          = make(map[string]*inventory_v1.Value)
		createdTime       = time.Now()
		updatedTime       = time.Now()

		randIntValue = gofakeit.Int64()

		modelPart = &model.Part{
			UUID:          partUUID,
			Name:          partName,
			Description:   partDescription,
			Price:         partPrice,
			StockQuantity: int64(partStockQuantity),
			Category:      model.CATEGORY_ENGINE,
			Dimensions:    &model.Dimensions{Length: length, Width: width, Height: height, Weight: weight},
			Manufacturer:  &model.Manufacturer{Name: manufacurerName, Country: manufacurerCountry, Website: manufacturerWebsite},
			Tags:          partTags,
			Metadata:      make(map[string]*model.Value),
			CreatedAt:     &createdTime,
			UpdatedAt:     &updatedTime,
		}

		invPart = &inventory_v1.Part{
			Uuid:          partUUID,
			Name:          partName,
			Description:   partDescription,
			Price:         partPrice,
			StockQuantity: int64(partStockQuantity),
			Category:      partCategory,
			Dimensions:    partDimensions,
			Manufacturer:  partManufacturer,
			Tags:          partTags,
			Metadata:      metadata,
			CreatedAt:     timestamppb.New(createdTime),
			UpdatedAt:     timestamppb.New(updatedTime),
		}
	)

	modelPart.Metadata[partUUID] = &model.Value{DataType: &model.Value_Int64Value{Int64Value: randIntValue}}
	modelPart.Metadata[partName] = &model.Value{DataType: &model.Value_Int64Value{Int64Value: randIntValue}}
	invPart.Metadata[partUUID] = &inventory_v1.Value{DataType: &inventory_v1.Value_Int64Value{Int64Value: randIntValue}}
	invPart.Metadata[partName] = &inventory_v1.Value{DataType: &inventory_v1.Value_Int64Value{Int64Value: randIntValue}}

	assert.Equal(t, ModelToProto(modelPart), invPart)
}

func TestRequestToModelPart_NilInput(t *testing.T) {
	assert.Nil(t, RequestToModelPart(nil))
}

func TestRequestToModelPart_Success(t *testing.T) {
	var (
		partUuids                 = []string{gofakeit.UUID(), gofakeit.UUID()}
		partNmaes                 = []string{gofakeit.ProductName(), gofakeit.ProductName()}
		partInventoryCategories   = []inventory_v1.Category{inventory_v1.Category_CATEGORY_ENGINE, inventory_v1.Category_CATEGORY_PORTHOLE}
		partManufacturerCountries = []string{gofakeit.Country(), gofakeit.Country()}
		partTags                  = []string{gofakeit.HackerPhrase(), gofakeit.HackerPhrase()}

		filter = &inventory_v1.PartsFilter{
			Uuids:                 partUuids,
			Names:                 partNmaes,
			Categories:            partInventoryCategories,
			ManufacturerCountries: partManufacturerCountries,
			Tags:                  partTags,
		}

		request = &inventory_v1.ListPartsRequest{
			Filter: filter,
		}

		partModelCategories = []model.Category{model.CATEGORY_ENGINE, model.CATEGORY_PORTHOLE}

		expectedPartsFilter = &model.PartsFilter{
			Uuids:                 partUuids,
			Names:                 partNmaes,
			Categories:            partModelCategories,
			ManufacturerCountries: partManufacturerCountries,
			Tags:                  partTags,
		}
	)
	assert.Equal(t, RequestToModelPart(request), expectedPartsFilter)
}

func TestModelListPartsToProto_NilInput(t *testing.T) {
	assert.Nil(t, ModelListPartsToProto(nil))
}

func TestModelListPartsToProto_Success(t *testing.T) {
	var (
		firstPartUUID  = gofakeit.UUID()
		secondPartUUID = gofakeit.UUID()
		parts          = []*model.Part{{UUID: firstPartUUID}, {UUID: secondPartUUID}}
		expectedParts  = []*inventory_v1.Part{{Uuid: firstPartUUID}, {Uuid: secondPartUUID}}
	)
	result := ModelListPartsToProto(parts)
	assert.Equal(t, result[0], expectedParts[0])
}
