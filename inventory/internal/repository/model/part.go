package inventory

import "time"

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Category int32

const (
	Category_CATEGORY_UNSPECIFIED Category = 0
	Category_CATEGORY_ENGINE      Category = 1
	Category_CATEGORY_FUEL        Category = 2
	Category_CATEGORY_PORTHOLE    Category = 3
	Category_CATEGORY_WING        Category = 4
)

// Enum value maps for Category.
var (
	Category_name = map[int32]string{
		0: "CATEGORY_UNSPECIFIED",
		1: "CATEGORY_ENGINE",
		2: "CATEGORY_FUEL",
		3: "CATEGORY_PORTHOLE",
		4: "CATEGORY_WING",
	}
	Category_value = map[string]int32{
		"CATEGORY_UNSPECIFIED": 0,
		"CATEGORY_ENGINE":      1,
		"CATEGORY_FUEL":        2,
		"CATEGORY_PORTHOLE":    3,
		"CATEGORY_WING":        4,
	}
)

type Dimensions struct {
	// Length in cm
	Length float64
	// Width in cm
	Width float64
	// Height in cm
	Height float64
	// Weight in kg
	Weight float64
}

type Manufacturer struct {
	// Name
	Name string
	// Country of origin
	Country string
	// Manufacturer's website
	Website string
}

type Value struct {
	// Types that are valid to be assigned to DataType:
	//
	//	*Value_StringValue
	//	*Value_Int64Value
	//	*Value_DoubleValue
	//	*Value_BoolValue
	DataType isValue_DataType
}

type isValue_DataType interface {
	isValue_DataType()
}

type Value_StringValue struct {
	// String value
	StringValue string
}

type Value_Int64Value struct {
	// Integer value
	Int64Value int64
}

type Value_DoubleValue struct {
	// Fractional value
	DoubleValue float64
}

type Value_BoolValue struct {
	// Logical value
	BoolValue bool
}

func (*Value_StringValue) isValue_DataType() {}

func (*Value_Int64Value) isValue_DataType() {}

func (*Value_DoubleValue) isValue_DataType() {}

func (*Value_BoolValue) isValue_DataType() {}
