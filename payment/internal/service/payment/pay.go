package payment

import (
	"context"
	"log"

	model "github.com/LushnikovSR/spaceship_factory/payment/internal/model"
	"github.com/google/uuid"
)

func (s *service) PayOrder(_ context.Context, orderID string, userID string, paymentMethod int32) (string, error) {
	// 1. Валидация обязательных полей
	if orderID == "" {
		return "", model.ErrMissingArgument
	}
	if userID == "" {
		return "", model.ErrMissingArgument
	}

	// 2. Проверка, что метод оплаты указан (не нулевое "неопределённое" значение)
	if paymentMethod == 0 { // предполагаем, что 0 = UNSPECIFIED
		return "", model.ErrMissingArgument
	}

	transactionUuid := uuid.NewString()

	log.Printf("Successful payment, transaction_uuid: %s", transactionUuid)

	return transactionUuid, nil
}
