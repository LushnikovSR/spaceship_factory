package inventory

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	data *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		data: db.Collection("parts"),
	}
}
