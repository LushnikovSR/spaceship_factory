package inventory

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	data *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	repo := &repository{
		data: db.Collection("parts"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	repo.Init(ctx)

	return repo
}
