package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	orderV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

const (
	httpPort          = "8080"
	readHeaderTimeout = 5 * time.Second
	contextTimeout    = 10 * time.Second
)

// OrderStorage представляет потокобезопасное хранилище данных о заказах
type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

// NewOrderStorage создает новое хранилище данных о заказах
func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

// GetOrder возвращает информацию о заказе по uuid.
// Если заказ не найден, возвращает nil.
func (s *OrderStorage) GetOrder(uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return order
}

// UpdateOrder обновляет данные о заказе для указанного заказа.
// Если заказа нет в хранилище, создает новую запись.
func (s *OrderStorage) UpdateOrder(order *orderV1.OrderDto) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[order.OrderUUID] = order
}

// DeleteOrder удаляет данные о заказе по указанному uuid.
// Метод ничего не возвращает, в том числе для случая когда заказа нет в хранилище или он уже был удалён.
func (s *OrderStorage) DeleteOrder(uuid string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.orders, uuid)
}

// OrderHandler реализует интерфейс orderV1.Handler для обработки запросов к API заказа
type OrderHandler struct {
	storage *OrderStorage
}

// NewOrderHandler создает новый обработчик запросов к API заказа
func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

// CancelOrder implements cancelOrder operation.
//
// Checks the order status. If `PENDING_PAYMENT`, changes the status to `CANCELLED`. If `PAID`,
// returns a 409 error.
//
// POST /orders/{order_uuid}/cancel
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order for uuid '" + params.OrderUUID + "' not found",
		}, nil
	}

	if order.Status == orderV1.OrderDtoStatusPAID {
		return &orderV1.ConflictError{
			Code:    409,
			Message: "The order has been paid. Cancellation is not possible.",
		}, nil
	}

	if order.Status == orderV1.OrderDtoStatusCANCELLED {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "The order has already been cancelled. Cancellation is not possible again.",
		}, nil
	}

	order.Status = orderV1.OrderDtoStatusPAID

	h.storage.UpdateOrder(order)

	return orderV1.CancelOrderNoContent, nil
}

// CreateOrder implements createOrder operation.
//
// Получает детали через `InventoryService.ListParts`. Проверяет, что
// все детали существуют. Если хотя бы одной нет —
// возвращает ошибку. Считает `total_price`. Генерирует `order_uuid`.
//  Сохраняет заказ со статусом `PENDING_PAYMENT`.
//
// POST /orders
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	parts, err := InventoryService.ListParts(req.PartUuids)
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: "The specified part is not in stock.",
		}, err
	}

	var total_price float64 = 0.0
	for _, part := range parts {
		total_price += part.Price
	}

	orderUUID := GetUniqueUUID(h.storage)

	order := orderV1.OrderDto{
		OrderUUID:       orderUUID,
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      total_price,
		TransactionUUID: orderV1.OptNilString.SetToNull(),
		PaymentMethod:   orderV1.NilOrderDtoPaymentMethod.SetToNull(),
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
	}

	h.storage.UpdateOrder(order)

	return req, nil
}

func GetUniqueUUID(storage *OrderStorage) string {
	var unique bool
	var randUuid string
	for {
		randUuid = uuid.NewString()

		order := storage.GetOrder(randUuid)
		if order == nil {
			unique = true
		}

		if unique {
			break
		}
	}

	return randUuid
}

// GetOrder implements getOrder operation.
//
// Get order by uuid.
//
// GET /orders/{order_uuid}
func (h *OrderHandler) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order for uuid '" + params.OrderUUID + "' not found",
		}, nil
	}

	return order, nil
}

// PayOrder implements payOrder operation.
//
// Находит заказ по `order_uuid`. Если не существует —
// возвращает 404 Not Found. Вызывает `PaymentService.PayOrder`, передаёт
// `user_uuid`, `order_uuid` и `payment_method`. Получает`transaction_uuid`.
// Обновляет заказ: статус → `PAID`, сохраняет `transaction_uuid`,
// `payment_method`.
//
// POST /orders/{order_uuid}/pay
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
	if order == nil {
		return &orderV1.NotFoundError{
			Code:    404,
			Message: "Order for uuid '" + params.OrderUUID + "' not found",
		}, nil
	}
	payMethod := req.GetPaymentMethod()

	transaction_uuid, err := PaymentService.PayOrder(order.UserUUID, order.OrderUUID, payMethod)
	if err != nil {
		return &orderV1.InternalError{
			Code:    500,
			Message: "PaymentService Error",
		}, nil
	}

	order.TransactionUUID.SetTo(transaction_uuid)
	order.PaymentMethod.SetTo(payMethod)
	order.SetStatus(orderV1.OrderDtoStatusPAID)

	h.storage.UpdateOrder(order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: transaction_uuid,
	}, nil
}

// Функция NewError создает *GenericErrorStatusCode из ошибки, возвращенной обработчиком.
// Используется для стандартного ответа по умолчанию.
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    orderV1.NewOptInt(http.StatusInternalServerError),
			Message: orderV1.NewOptString(err.Error()),
		},
	}
}

func main() {

	//Создаём новое хранилище для данных о заказах
	storage := NewOrderStorage()

	//Создаём обработчик API заказов
	orderHandler := NewOrderHandler(storage)

	//Создаём ОpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	// Инициализируем роутер Chi
	r := chi.NewRouter()

	//Добавляем middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Монтируем обработчики OpenAPI
	r.Mount("/", orderServer)

	//Запускаем http-server
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, //для защиты от slowloris атак
	}

	//Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server is starting on port: %s\n", httpPort)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Starting server is failed: %s\n", err)
		}
	}()

	//Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server stopping ...")

	//Создаём контекст с таймутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	err = orderServer.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Error during server closing: %v\n", err)
	}
	log.Printf("Server is stoped correctly")
}
