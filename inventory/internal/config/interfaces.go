package inventory

type LoggerConfig interface {
	Level() string
	AsJson() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
