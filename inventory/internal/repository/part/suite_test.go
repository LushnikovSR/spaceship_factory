package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	ctx        context.Context
	repository *repository
}

func (s *RepositorySuite) SetupTest() {
	s.ctx = context.Background()
	s.repository = NewRepository()
}

func (s *RepositorySuite) TearDownTest() {}

func TestRepositoryIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
