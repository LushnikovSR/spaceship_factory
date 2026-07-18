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

func TestRepository_ListParts(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		repo := NewRepository(mt.DB)

		repoPart := repoModel.Part{
			ID:            primitive.NewObjectID(),
			Name:          "Engine",
			Price:         1000,
			StockQuantity: 5,
			Category:      repoModel.CATEGORY_ENGINE,
			Tags:          []string{"engine"},
		}

		raw, err := bson.Marshal(repoPart)
		require.NoError(mt, err)

		var doc bson.D
		require.NoError(mt, bson.Unmarshal(raw, &doc))

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"inventory.parts",
				mtest.FirstBatch,
				doc,
			),
			mtest.CreateCursorResponse(
				0,
				"inventory.parts",
				mtest.NextBatch,
			),
		)

		parts, err := repo.ListParts(context.Background(), nil)

		require.NoError(mt, err)
		require.Len(mt, parts, 1)

		require.Equal(mt, repoPart.ID.Hex(), parts[0].UUID)
		require.Equal(mt, repoPart.Name, parts[0].Name)
		require.EqualValues(mt, repoPart.StockQuantity, parts[0].StockQuantity)
	})

	mt.Run("empty", func(mt *mtest.T) {
		repo := NewRepository(mt.DB)

		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				0,
				"inventory.parts",
				mtest.FirstBatch,
			),
		)

		parts, err := repo.ListParts(context.Background(), nil)

		require.NoError(mt, err)
		require.Empty(mt, parts)
	})
}

func TestBuildMongoFilter(t *testing.T) {
	engineID := primitive.NewObjectID()

	tests := []struct {
		name    string
		filter  *model.PartsFilter
		want    bson.M
		wantErr bool
	}{
		{
			name:   "nil filter",
			filter: nil,
			want:   bson.M{},
		},
		{
			name: "name filter",
			filter: &model.PartsFilter{
				Names: []string{"Engine"},
			},
			want: bson.M{
				"name": bson.M{
					"$in": []string{"Engine"},
				},
			},
		},
		{
			name: "manufacturer country",
			filter: &model.PartsFilter{
				ManufacturerCountries: []string{"Germany"},
			},
			want: bson.M{
				"manufacturer.country": bson.M{
					"$in": []string{"Germany"},
				},
			},
		},
		{
			name: "tags",
			filter: &model.PartsFilter{
				Tags: []string{"engine", "main"},
			},
			want: bson.M{
				"tags": bson.M{
					"$in": []string{"engine", "main"},
				},
			},
		},
		{
			name: "uuid",
			filter: &model.PartsFilter{
				Uuids: []string{engineID.Hex()},
			},
			want: bson.M{
				"_id": bson.M{
					"$in": []primitive.ObjectID{engineID},
				},
			},
		},
		{
			name: "invalid uuid",
			filter: &model.PartsFilter{
				Uuids: []string{"invalid"},
			},
			wantErr: true,
		},
		{
			name: "combined",
			filter: &model.PartsFilter{
				Uuids:                 []string{engineID.Hex()},
				Names:                 []string{"Engine"},
				ManufacturerCountries: []string{"Germany"},
				Tags:                  []string{"engine"},
			},
			want: bson.M{
				"_id": bson.M{
					"$in": []primitive.ObjectID{engineID},
				},
				"name": bson.M{
					"$in": []string{"Engine"},
				},
				"manufacturer.country": bson.M{
					"$in": []string{"Germany"},
				},
				"tags": bson.M{
					"$in": []string{"engine"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildMongoFilter(tt.filter)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
