package order

import (
	"context"
	"errors"
	"fmt"

	model "github.com/LushnikovSR/spaceship_factory/order/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/order/internal/repository/converter"
	sq "github.com/Masterminds/squirrel"
)

// UpdateOrder обновляет данные о заказе для указанного заказа.
// Если заказа нет в хранилище, создает новую запись.
func (r *repository) UpdateOrder(ctx context.Context, order *model.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}
	// Преобразуем в модель репозитория для корректных типов
	ro := converter.OrderModelToRepoModel(order)

	// Строим запрос на обновление записи в таблице orders
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("transaction_uuid", ro.TransactionUUID).
		Set("payment_method", ro.PaymentMethod).
		Set("status", ro.Status).
		Where(sq.Eq{"order_uuid": ro.OrderUUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	res, err := r.data.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}
	fmt.Printf("\nfor id %s updated %d rows\n", order.OrderUUID, res.RowsAffected())

	return nil
}
