package testcontainers

const (
	// MongoDB conteiner constants
	MongoConteinerName = "mongo"
	MongoPort          = "27017"

	// MongoDN environment constants
	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD"
	MongoAuthKey      = "MONGO_AUTH_DB"
)
