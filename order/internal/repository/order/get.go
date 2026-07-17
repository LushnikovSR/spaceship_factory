package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/order/internal/repository/converter"
	repoModel "github.com/LushnikovSR/spaceship_factory/order/internal/repository/model"
)

// GetOrder возвращает информацию о заказе по uuid.
// Если заказ не найден, возвращает пустой объект.
func (r *repository) GetOrder(ctx context.Context, uuid string) (*model.Order, error) {
	// Строим запрос на выборку записей из таблицы orders
	builderSelect := sq.Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_uuid": uuid})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	ro := &repoModel.Order{}
	err = r.data.QueryRow(ctx, query, args...).Scan(&ro.OrderUUID, &ro.UserUUID, &ro.PartUuids, &ro.TotalPrice, &ro.TransactionUUID, &ro.PaymentMethod, &ro.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // заказ не найден — без ошибки
		}
		return nil, fmt.Errorf("failed to scan order: %w", err)
	}

	return converter.OrderRepoModelToModel(ro), nil
}
