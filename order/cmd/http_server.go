package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	customMiddleware "github.com/LushnikovSR/spaceship_factory/internal/middleware"
	orderV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort             = "8080"
	readHeaderTimeout    = 5 * time.Second
	contextTimeout       = 10 * time.Second
	inventoryServicePort = "50051"
	paymentServicePort   = "50052"
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
	storage         *OrderStorage
	inventoryClient inventory_v1.InventoryServiceClient
	paymentClient   payment_v1.PaymentServiceClient
}

// NewOrderHandler создает новый обработчик запросов к API заказа
func NewOrderHandler(storage *OrderStorage,
	invClient inventory_v1.InventoryServiceClient,
	payClient payment_v1.PaymentServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		inventoryClient: invClient,
		paymentClient:   payClient,
	}
}

// CancelOrder implements cancelOrder operation.
//
// Checks the order status. If `PENDING_PAYMENT`, changes the status to `CANCELLED`. If `PAID`,
// returns a 409 error.
//
// POST /orders/{order_uuid}/cancel
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order := h.storage.GetOrder(params.OrderUUID)
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

	order.SetStatus(orderV1.OrderDtoStatusCANCELLED)
	h.storage.UpdateOrder(order)

	return &orderV1.CancelOrderNoContent{}, nil
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
	// Получаем детали из InventoryService
	listResp, err := h.inventoryClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: req.PartUuids,
		},
	})
	if err != nil {
		// Ошибка связи с InventoryService
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Failed to fetch parts from inventory",
		}, err
	}

	// Проверяем, что все запрошенные детали существуют
	foundUuids := make(map[string]struct{}, len(listResp.Parts))
	for _, p := range listResp.Parts {
		foundUuids[p.Uuid] = struct{}{}
	}
	for _, uid := range req.PartUuids {
		if _, ok := foundUuids[uid]; !ok {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Part with UUID " + uid + " not found",
			}, err
		}
	}

	// Считаем total_price
	var total_price float64 = 0.0
	for _, part := range listResp.Parts {
		total_price += part.Price
	}

	// Генерируем order_uuid и сохраняем заказ
	orderUUID := GetUniqueUUID(h.storage)

	var transactionUUID orderV1.OptNilString
	transactionUUID.SetToNull()

	var paymentMethod orderV1.NilOrderDtoPaymentMethod
	paymentMethod.SetToNull()

	order := orderV1.OrderDto{
		OrderUUID:       orderUUID,
		UserUUID:        req.UserUUID,
		PartUuids:       req.PartUuids,
		TotalPrice:      total_price,
		TransactionUUID: transactionUUID,
		PaymentMethod:   &paymentMethod,
		Status:          orderV1.OrderDtoStatusPENDINGPAYMENT,
	}

	h.storage.UpdateOrder(&order)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: total_price,
	}, nil
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

	paymentMethod := orderV1.OrderDtoPaymentMethod(req.GetPaymentMethod())

	payResp, err := h.paymentClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     order.OrderUUID,
		UserUuid:      order.UserUUID,
		PaymentMethod: mapPaymentMethod(paymentMethod),
	})
	if err != nil {
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Payment service error: " + err.Error(),
		}, nil
	}

	order.TransactionUUID.SetTo(payResp.TransactionUuid)
	order.PaymentMethod.SetTo(paymentMethod)
	order.SetStatus(orderV1.OrderDtoStatusPAID)

	h.storage.UpdateOrder(order)

	return &orderV1.PayOrderResponse{
		TransactionUUID: payResp.TransactionUuid,
	}, nil
}

// Функция NewError создает *GenericErrorStatusCode из ошибки, возвращенной обработчиком.
// Используется для стандартного ответа по умолчанию.
func (h *OrderHandler) NewError(_ context.Context, err error) *orderV1.GenericErrorStatusCode {
	return &orderV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: orderV1.GenericError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		},
	}
}

// Вспомогательная функция для конвертации метода оплаты
func mapPaymentMethod(m orderV1.OrderDtoPaymentMethod) payment_v1.PaymentMethod {
	switch m {
	case orderV1.OrderDtoPaymentMethodCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.OrderDtoPaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.OrderDtoPaymentMethodCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.OrderDtoPaymentMethodINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func main() {
	// gRPC подключение к InventoryService
	connInv, err := grpc.NewClient(fmt.Sprintf("localhost:%s", inventoryServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to InventoryService: %v", err)
	}
	defer func() {
		if cerr := connInv.Close(); cerr != nil {
			log.Printf("failed to close InventoryService connection: %v", cerr)
		}
	}()
	invClient := inventory_v1.NewInventoryServiceClient(connInv)

	//gRPC подключение к PaymentService
	connPay, err := grpc.NewClient(fmt.Sprintf("localhost:%s", paymentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to PaymentService: %v", err)
	}
	defer func() {
		if cerr := connPay.Close(); cerr != nil {
			log.Printf("failed to close PaymentService connection: %v", cerr)
		}
	}()
	payClient := payment_v1.NewPaymentServiceClient(connPay)

	//Создаём новое хранилище для данных о заказах
	storage := NewOrderStorage()

	//Создаём обработчик API заказов
	orderHandler := NewOrderHandler(storage, invClient, payClient)

	//Создаём ОpenAPI сервер
	orderServer, err := orderV1.NewServer(orderHandler, orderV1.WithPathPrefix("/api/v1"))
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
	r.Use(customMiddleware.RequestLogger)

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

	err = server.Shutdown(ctx)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Error during server closing: %v\n", err)
	}
	log.Printf("Server is stoped correctly")
}
