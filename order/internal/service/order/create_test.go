package order

import (
	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateOrder_Success() {
	var (
		userUUID  = gofakeit.UUID()
		partUuids = []string{
			"11111111-1111-1111-1111-111111111111",
			"22222222-2222-2222-2222-222222222222",
		}
		expectedTotal = 1800.0

		parts = []model.Part{
			{
				Uuid:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Tags:          []string{"engine", "main"},
			},
			{
				Uuid:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Tags:          []string{"porthole", "window"},
			},
		}
	)

	// 1. Мокаем ListParts – сервис вызовет его с PartsFilter{Uuids: partUuids}
	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{Uuids: partUuids}).
		Return(parts, nil).
		Once()

	// 2. Мокаем CreateOrder – принимает любой *model.Order, возвращает nil
	s.orderRepository.
		On("CreateOrder", mock.MatchedBy(func(o any) bool {
			// Приведение к model.Order (или к тому типу, который реально используется)
			order, ok := o.(*model.Order)
			return ok && order.UserUUID == userUUID && len(order.PartUuids) == len(partUuids)
		})).
		Return(nil).
		Once()

	s.orderRepository.
		On("GetOrder", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	// 3. Вызываем тестируемый метод
	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	// 4. Проверяем результаты
	s.Require().NoError(err)
	s.Require().NotEmpty(orderUUID)
	s.Require().Len(orderUUID, 36)
	s.Require().Equal(expectedTotal, totalPrice)
}

func (s *ServiceSuite) TestCreateOrder_InventoryError() {
	userUUID := gofakeit.UUID()
	partUuids := []string{"any-uuid"}

	// ListParts возвращает ошибку
	s.inventoryClient.
		On("ListParts", s.ctx, mock.Anything).
		Return(nil, gofakeit.Error()).
		Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().Error(err)
	// Проверяем, что это InternalServerError (код 500)
	var internalErr *model.InternalServerError
	s.Require().ErrorAs(err, &internalErr)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
}

func (s *ServiceSuite) TestCreateOrder_PartNotFound() {
	userUUID := gofakeit.UUID()
	partUuids := []string{"11111111-1111-1111-1111-111111111111", "non-existent"}

	parts := []model.Part{
		{Uuid: "11111111-1111-1111-1111-111111111111", Price: 100},
	}

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{Uuids: partUuids}).
		Return(parts, nil).
		Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().Error(err)
	var notFoundErr *model.NotFoundError
	s.Require().ErrorAs(err, &notFoundErr)
	s.Require().Contains(err.Error(), "non-existent")
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
}

func (s *ServiceSuite) TestCreateOrder_RepositoryError() {
	userUUID := gofakeit.UUID()
	partUuids := []string{"11111111-1111-1111-1111-111111111111"}

	parts := []model.Part{
		{Uuid: "11111111-1111-1111-1111-111111111111", Price: 500},
	}

	s.inventoryClient.
		On("ListParts", s.ctx, mock.Anything).
		Return(parts, nil).
		Once()

	s.orderRepository.
		On("CreateOrder", mock.MatchedBy(func(o any) bool {
			// Приведение к model.Order (или к тому типу, который реально используется)
			order, ok := o.(*model.Order)
			return ok && order.UserUUID == userUUID && len(order.PartUuids) == len(partUuids)
		})).
		Return(gofakeit.Error()).
		Once()

	s.orderRepository.
		On("GetOrder", mock.AnythingOfType("string")).
		Return(nil).
		Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().Error(err)
	var conflictErr *model.ConflictError
	s.Require().ErrorAs(err, &conflictErr)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
}
