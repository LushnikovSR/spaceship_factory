package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type RepositorySuite struct {
	suite.Suite
	ctx        context.Context
	repository *repository
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()
	// Создаём mock-клиент MongoDB (in-memory) для каждого теста
	mt := mtest.New(s.T(), mtest.NewOptions().ClientType(mtest.Mock))
	db := mt.DB // *mongo.Database
	s.repository = NewRepository(db)
}

func (s *RepositorySuite) TearDownTest() {
	// Очищаем коллекцию после теста для изоляции
	_, _ = s.repository.data.DeleteMany(s.ctx, bson.D{})
}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
