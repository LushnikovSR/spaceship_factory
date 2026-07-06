package inventory

import (
	"testing"
	"time"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRepoModelToModel_NilInput(t *testing.T) {
	assert.Nil(t, RepoModelToModel(nil))
}

func TestRepoModelToModel_FullDataSuccess(t *testing.T) {
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
		partCategory      = repoModel.Category_CATEGORY_ENGINE
		partDimensions    = &repoModel.Dimensions{Length: length, Width: width, Height: height, Weight: weight}
		partManufacturer  = &repoModel.Manufacturer{Name: manufacurerName, Country: manufacurerCountry, Website: manufacturerWebsite}
		partTags          = []string{}
		createdTime       = time.Now()
		updatedTime       = time.Now()

		randIntValue = gofakeit.Int64()
		randFloat    = gofakeit.Float64Range(1, 1000)
		randBool     = gofakeit.Bool()

		metadata      = make(map[string]*repoModel.Value)
		modelMetadata = make(map[string]*model.Value)

		expectedPart = &model.Part{
			UUID:          partUUID,
			Name:          partName,
			Description:   partDescription,
			Price:         partPrice,
			StockQuantity: int64(partStockQuantity),
			Category:      model.Category_CATEGORY_ENGINE,
			Dimensions:    &model.Dimensions{Length: length, Width: width, Height: height, Weight: weight},
			Manufacturer:  &model.Manufacturer{Name: manufacurerName, Country: manufacurerCountry, Website: manufacturerWebsite},
			Tags:          partTags,
			Metadata:      modelMetadata,
			CreatedAt:     &createdTime,
			UpdatedAt:     &updatedTime,
		}

		part = &repoModel.Part{
			UUID:          partUUID,
			Name:          partName,
			Description:   partDescription,
			Price:         partPrice,
			StockQuantity: int64(partStockQuantity),
			Category:      partCategory,
			Dimensions:    partDimensions,
			Manufacturer:  partManufacturer,
			Tags:          partTags,
			Metadata:      metadata,
			CreatedAt:     &createdTime,
			UpdatedAt:     &updatedTime,
		}

		key1  = gofakeit.UUID()
		key2  = gofakeit.UUID()
		key3  = gofakeit.UUID()
		key4  = gofakeit.UUID()
		key5  = gofakeit.UUID()
		key6  = gofakeit.UUID()
		key7  = gofakeit.UUID()
		key8  = gofakeit.UUID()
		key9  = gofakeit.UUID()
		key10 = gofakeit.UUID()
	)

	part.Metadata[key1] = &repoModel.Value{DataType: &repoModel.Value_Int64Value{Int64Value: randIntValue}}
	part.Metadata[key2] = &repoModel.Value{DataType: &repoModel.Value_Int64Value{Int64Value: randIntValue}}
	part.Metadata[key3] = &repoModel.Value{DataType: &repoModel.Value_StringValue{StringValue: partUUID}}
	part.Metadata[key4] = &repoModel.Value{DataType: &repoModel.Value_StringValue{StringValue: partName}}
	part.Metadata[key5] = &repoModel.Value{DataType: &repoModel.Value_DoubleValue{DoubleValue: randFloat}}
	part.Metadata[key6] = &repoModel.Value{DataType: &repoModel.Value_DoubleValue{DoubleValue: randFloat}}
	part.Metadata[key7] = &repoModel.Value{DataType: &repoModel.Value_BoolValue{BoolValue: randBool}}
	part.Metadata[key8] = &repoModel.Value{DataType: &repoModel.Value_BoolValue{BoolValue: randBool}}
	part.Metadata[key9] = nil
	part.Metadata[key10] = nil

	expectedPart.Metadata[key1] = &model.Value{DataType: &model.Value_Int64Value{Int64Value: randIntValue}}
	expectedPart.Metadata[key2] = &model.Value{DataType: &model.Value_Int64Value{Int64Value: randIntValue}}
	expectedPart.Metadata[key3] = &model.Value{DataType: &model.Value_StringValue{StringValue: partUUID}}
	expectedPart.Metadata[key4] = &model.Value{DataType: &model.Value_StringValue{StringValue: partName}}
	expectedPart.Metadata[key5] = &model.Value{DataType: &model.Value_DoubleValue{DoubleValue: randFloat}}
	expectedPart.Metadata[key6] = &model.Value{DataType: &model.Value_DoubleValue{DoubleValue: randFloat}}
	expectedPart.Metadata[key7] = &model.Value{DataType: &model.Value_BoolValue{BoolValue: randBool}}
	expectedPart.Metadata[key8] = &model.Value{DataType: &model.Value_BoolValue{BoolValue: randBool}}
	expectedPart.Metadata[key9] = nil
	expectedPart.Metadata[key10] = nil

	assert.Equal(t, RepoModelToModel(part), expectedPart)
}

func TestRepoModelToModel_OnlyUUIDSuccess(t *testing.T) {
	var (
		partUUID = gofakeit.UUID()

		expectedPart = &model.Part{
			UUID: partUUID,
		}

		part = &repoModel.Part{
			UUID: partUUID,
		}
	)

	assert.Equal(t, RepoModelToModel(part), expectedPart)
}
