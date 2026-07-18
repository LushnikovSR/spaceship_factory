package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
)

func TestRepository_GetPart(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	ctx := context.Background()

	mt.Run("success", func(mt *mtest.T) {
		repo := NewRepository(mt.DB)

		id := primitive.NewObjectID()

		repoPart := repoModel.Part{
			ID:            id,
			Name:          "Сопло маршевое",
			Price:         1500,
			StockQuantity: 5,
			Category:      repoModel.CATEGORY_ENGINE,
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Biscuit",
				Country: "Germany",
				Website: "financialharness.info",
			},
			Tags: []string{"engine", "main"},
		}

		raw, err := bson.Marshal(repoPart)
		require.NoError(mt, err)

		var doc bson.D
		require.NoError(mt, bson.Unmarshal(raw, &doc))

		mt.AddMockResponses(
			mtest.CreateCursorResponse(1, "inventory.parts", mtest.FirstBatch, doc),
			mtest.CreateCursorResponse(0, "inventory.parts", mtest.NextBatch),
		)

		part, err := repo.GetPart(ctx, id.Hex())

		require.NoError(mt, err)
		require.NotNil(mt, part)

		require.Equal(mt, id.Hex(), part.UUID)
		require.Equal(mt, "Сопло маршевое", part.Name)
		require.Equal(mt, 1500.0, part.Price)
		require.EqualValues(mt, 5, part.StockQuantity)
		require.Equal(mt, model.CATEGORY_ENGINE, part.Category)

		require.NotNil(mt, part.Manufacturer)
		require.Equal(mt, "Biscuit", part.Manufacturer.Name)
		require.Equal(mt, "Germany", part.Manufacturer.Country)
		require.Equal(mt, "financialharness.info", part.Manufacturer.Website)

		require.Equal(mt, []string{"engine", "main"}, part.Tags)
	})

	mt.Run("not found", func(mt *mtest.T) {
		repo := NewRepository(mt.DB)

		id := primitive.NewObjectID()

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				0,
				"inventory.parts",
				mtest.FirstBatch,
			),
		)

		part, err := repo.GetPart(ctx, id.Hex())

		require.Error(mt, err)
		require.ErrorIs(mt, err, model.ErrPartNotFound)
		require.NotNil(mt, part)
		require.Equal(mt, &model.Part{}, part)
	})

	mt.Run("invalid uuid", func(mt *mtest.T) {
		repo := NewRepository(mt.DB)

		part, err := repo.GetPart(ctx, "invalid-id")

		require.Error(mt, err)
		require.Empty(mt, part)
	})
}
