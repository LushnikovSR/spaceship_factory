package order

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	grpcMocks "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/mocks"
	"github.com/LushnikovSR/spaceship_factory/order/internal/repository/mocks"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/logger"
)

var loggerLevelValue = "debug"

type ServiceSuite struct {
	suite.Suite
	ctx             context.Context
	orderRepository *mocks.OrderRepository
	inventoryClient *grpcMocks.InventoryClient
	paymentClient   *grpcMocks.PaymentClient
	service         *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.orderRepository = mocks.NewOrderRepository(s.T())
	s.inventoryClient = grpcMocks.NewInventoryClient(s.T())
	s.paymentClient = grpcMocks.NewPaymentClient(s.T())
	s.service = NewService(s.orderRepository, s.inventoryClient, s.paymentClient)

	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Errorf("не удалось инициализировать логгер: %w", err))
	}
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
