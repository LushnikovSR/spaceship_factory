package order

import (
	"context"
	"fmt"
	"reflect"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	logger "github.com/LushnikovSR/spaceship_factory/platform/pkg/logger"
	"go.uber.org/zap"
)

// GetOrder implements getOrder operation.
//
// Get order by uuid.
//
// GET /orders/{order_uuid}
func (s *service) GetOrder(ctx context.Context, orderUUID string) (model.Order, error) {
	order, err := s.orderRepository.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Error(ctx, "failed to get order", zap.String("orderUUID", orderUUID), zap.Error(err))
		return model.Order{}, &model.InternalServerError{
			BaseError: model.BaseError{
				Code:    500,
				Message: fmt.Errorf("order for uuid %v not found: %w", orderUUID, err).Error(),
			},
		}
	}

	empty := &model.Order{}

	if reflect.DeepEqual(order, empty) || order == nil {
		logger.Error(ctx, "order was not found", zap.String("orderUUID", orderUUID), zap.Error(err))
		return model.Order{}, &model.NotFoundError{
			BaseError: model.BaseError{
				Code:    404,
				Message: "order for uuid '" + orderUUID + "' not found",
			},
		}
	}

	return *order, nil
}
