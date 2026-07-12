package inventory

import (
	"log/slog"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
)

func TestRepoModelToModel_EmptyInput(t *testing.T) {
	assert.Empty(t, RepoModelToModel(repoModel.Part{}))
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

		randID            = gofakeit.Regex("[0-9a-f]{24}")
		partName          = gofakeit.Name()
		partDescription   = gofakeit.HackerPhrase()
		partPrice         = gofakeit.Price(100.0, 2000.0)
		partStockQuantity = gofakeit.RandomInt([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
		partCategory      = repoModel.CATEGORY_ENGINE
		partDimensions    = &repoModel.Dimensions{Length: length, Width: width, Height: height, Weight: weight}
		partManufacturer  = &repoModel.Manufacturer{Name: manufacurerName, Country: manufacurerCountry, Website: manufacturerWebsite}
		partTags          = []string{}
		createdTime       = time.Now()
		updatedTime       = time.Now()

		randIntValue = gofakeit.Int64()
		randFloat    = gofakeit.Float64Range(1, 1000)
		randBool     = gofakeit.Bool()

		metadata      = make(map[string]any)
		modelMetadata = make(map[string]*model.Value)

		expectedPart = &model.Part{
			Name:          partName,
			Description:   partDescription,
			Price:         partPrice,
			StockQuantity: int64(partStockQuantity),
			Category:      model.CATEGORY_ENGINE,
			Dimensions:    &model.Dimensions{Length: length, Width: width, Height: height, Weight: weight},
			Manufacturer:  &model.Manufacturer{Name: manufacurerName, Country: manufacurerCountry, Website: manufacturerWebsite},
			Tags:          partTags,
			Metadata:      modelMetadata,
			CreatedAt:     &createdTime,
			UpdatedAt:     &updatedTime,
		}

		part = repoModel.Part{
			Name:          partName,
			Description:   partDescription,
			Price:         partPrice,
			StockQuantity: int64(partStockQuantity),
			Category:      partCategory,
			Dimensions:    partDimensions,
			Manufacturer:  partManufacturer,
			Tags:          partTags,
			Metadata:      metadata,
			CreatedAt:     createdTime,
			UpdatedAt:     updatedTime,
		}

		key1  = gofakeit.Regex("[0-9a-f]{24}")
		key2  = gofakeit.Regex("[0-9a-f]{24}")
		key3  = gofakeit.Regex("[0-9a-f]{24}")
		key4  = gofakeit.Regex("[0-9a-f]{24}")
		key5  = gofakeit.Regex("[0-9a-f]{24}")
		key6  = gofakeit.Regex("[0-9a-f]{24}")
		key7  = gofakeit.Regex("[0-9a-f]{24}")
		key8  = gofakeit.Regex("[0-9a-f]{24}")
		key9  = gofakeit.Regex("[0-9a-f]{24}")
		key10 = gofakeit.Regex("[0-9a-f]{24}")
	)

	partID, err := primitive.ObjectIDFromHex(randID)
	if err != nil {
		slog.Debug("failed randID type string to type primitive.ObjectID")
	}
	part.ID = partID

	part.Metadata[key1] = randIntValue
	part.Metadata[key2] = randIntValue
	part.Metadata[key3] = randID
	part.Metadata[key4] = partName
	part.Metadata[key5] = randFloat
	part.Metadata[key6] = randFloat
	part.Metadata[key7] = randBool
	part.Metadata[key8] = randBool
	part.Metadata[key9] = nil
	part.Metadata[key10] = nil

	expectedPart.UUID = partID.Hex()

	expectedPart.Metadata[key1] = &model.Value{DataType: &model.Value_Int64Value{Int64Value: randIntValue}}
	expectedPart.Metadata[key2] = &model.Value{DataType: &model.Value_Int64Value{Int64Value: randIntValue}}
	expectedPart.Metadata[key3] = &model.Value{DataType: &model.Value_StringValue{StringValue: randID}}
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
		randID = gofakeit.Regex("[0-9a-f]{24}")

		expectedPart = &model.Part{}

		part = repoModel.Part{}
	)
	partID, err := primitive.ObjectIDFromHex(randID)
	if err != nil {
		slog.Debug("failed randID type string to type primitive.ObjectID")
	}
	part.ID = partID

	expectedPart.UUID = partID.Hex()

	assert.Equal(t, RepoModelToModel(part), expectedPart)
}
