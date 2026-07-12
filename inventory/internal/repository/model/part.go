package inventory

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Part – документ в коллекции parts.
type Part struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty"` // UUID детали (первичный ключ)
	Name          string                 `bson:"name"`
	Description   string                 `bson:"description,omitempty"`
	Price         float64                `bson:"price"`
	StockQuantity int64                  `bson:"stock_quantity"`
	Category      Category               `bson:"category"`
	Dimensions    *Dimensions            `bson:"dimensions,omitempty"`
	Manufacturer  *Manufacturer          `bson:"manufacturer,omitempty"`
	Tags          []string               `bson:"tags,omitempty"`
	Metadata      map[string]interface{} `bson:"metadata,omitempty"` // произвольные данные
	CreatedAt     time.Time              `bson:"created_at,omitempty"`
	UpdatedAt     time.Time              `bson:"updated_at,omitempty"`
}

// Category – тип категории детали.
type Category int32

const (
	CATEGORY_UNSPECIFIED Category = 0
	CATEGORY_ENGINE      Category = 1
	CATEGORY_FUEL        Category = 2
	CATEGORY_PORTHOLE    Category = 3
	CATEGORY_WING        Category = 4
)

// Dimensions – габариты и вес детали.
type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

// Manufacturer – производитель детали.
type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}
