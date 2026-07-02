package order

import "time"

//PartsFilter contains fields for filtering.
//If field is empty means no filtering by field.
//(All fields are optional)
type PartsFilter struct {
	//List of UUIDs (optional)
	Uuids []string
	//List of names (optional)
	Names []string
	//List of categories (optional)
	Categories []Category
	//List of manufacturing countries (optional)
	ManufacturerCountries []string
	//List of tags (optional)
	Tags []string
}

//Part provides complete information about the part
type Part struct {
	//uuid unique part identifier
	Uuid string
	//name of part
	Name string
	//description of part
	Description string
	//unit price
	Price float64
	//quantity in stock
	StockQuantity int64
	//category
	Category Category
	//dimensions of the part
	Dimensions *Dimensions
	//information about the manufacturer
	Manufacturer *Manufacturer
	//tags for quick search
	Tags []string
	//flexible metadata
	Metadata map[string]*Value
	//date of creation
	CreatedAt *time.Time
	//date of updation
	UpdatedAt *time.Time
}

//Category contains a list of possible categories
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

//Dimensions provides complete information about the dimension
type Dimensions struct {
	//Length in cm
	Length float64
	//Width in cm
	Width float64
	//Height in cm
	Height float64
	//Weight in kg
	Weight float64
}

//Manufacturer provides complete information about the manufacturer
type Manufacturer struct {
	//Name
	Name string
	//Country of origin
	Country string
	//Manufacturer's website
	Website string
}

//Value provides complete information about the
type Value struct {
	// Types that are valid to be assigned to DataType:
	//
	//	*Value_StringValue
	//	*Value_Int64Value
	//	*Value_DoubleValue
	//	*Value_BoolValue
	DataType isValue_DataType
}

func (x *Value) GetInt64Value() int64 {
	if x != nil {
		if x, ok := x.DataType.(*Value_Int64Value); ok {
			return x.Int64Value
		}
	}
	return 0
}

func (x *Value) GetDoubleValue() float64 {
	if x != nil {
		if x, ok := x.DataType.(*Value_DoubleValue); ok {
			return x.DoubleValue
		}
	}
	return 0
}

func (x *Value) GetBoolValue() bool {
	if x != nil {
		if x, ok := x.DataType.(*Value_BoolValue); ok {
			return x.BoolValue
		}
	}
	return false
}

type isValue_DataType interface {
	isValue_DataType()
}

type Value_StringValue struct {
	//String value
	StringValue string
}

type Value_Int64Value struct {
	//Integer value
	Int64Value int64
}

type Value_DoubleValue struct {
	//Fractional value
	DoubleValue float64
}

type Value_BoolValue struct {
	//Logical value
	BoolValue bool
}

func (*Value_StringValue) isValue_DataType() {}

func (*Value_Int64Value) isValue_DataType() {}

func (*Value_DoubleValue) isValue_DataType() {}

func (*Value_BoolValue) isValue_DataType() {}
