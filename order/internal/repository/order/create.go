package order

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/order/internal/repository/converter"
)

// CreateOrder создает новую запись о заказе.
func (r *repository) CreateOrder(ctx context.Context, order *model.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	// Преобразуем в модель репозитория
	ro := converter.OrderModelToRepoModel(order)

	// Строим запрос на вставку записи в таблицу orders
	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").
		Values(ro.UserUUID, ro.PartUuids, ro.TotalPrice, ro.TransactionUUID, ro.PaymentMethod, ro.Status).
		Suffix("RETURNING order_uuid")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	var orderID string
	err = r.data.QueryRow(ctx, query, args...).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
	}

	order.OrderUUID = orderID

	return nil
}
