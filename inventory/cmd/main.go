package inventory

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"

	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	grpcPort = 50051
)

// PartStorage представляет потокобезопасное хранилище данных о деталях
type PartStorage struct {
	mu      sync.RWMutex
	storage map[string]*inventory_v1.Part
}

func NewPartStorage() *PartStorage {
	return &PartStorage{
		storage: make(map[string]*inventory_v1.Part),
	}
}

// GetPart возвращает информацию о запчасте по uuid.
// Если запчасть не найдена, возвращает nil.
func (s *PartStorage) GetPart(uuid string) *inventory_v1.Part {
	s.mu.RLock()
	defer s.mu.RUnlock()

	part, ok := s.storage[uuid]
	if !ok {
		return nil
	}
	return part
}

// GetParts возвращает информацию о запчастях по списку uuid.
// Если запчасть не найдена, возвращает nil.
func (s *PartStorage) GetParts(uuids []string) []*inventory_v1.Part {
	seen := make(map[string]struct{}, len(uuids))
	parts := make([]*inventory_v1.Part, 0, len(uuids))

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, id := range uuids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if part, ok := s.storage[id]; ok {
			parts = append(parts, part)
		}
	}

	return parts
}

// GetAllParts возвращает информацию о всех запчастях в хранилище.
func (s *PartStorage) GetAllParts() []*inventory_v1.Part {
	parts := make([]*inventory_v1.Part, 0)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, part := range s.storage {
		parts = append(parts, part)
	}

	return parts
}

// UpdateStorage обновляет данные о запчасте.
// Если запчасти нет в хранилище, создает новую запись.
func (s *PartStorage) UpdatePart(part *inventory_v1.Part) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[part.Uuid] = part
}

// DeletePart удаляет данные о запчасти по указанному uuid.
// Метод ничего не возвращает, в том числе для случая когда запчасти нет в хранилище или она уже была удалёна.
func (s *PartStorage) DeletePart(uuid string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.storage, uuid)
}

// InventoryService реализует gRPC сервис для работы с хранилищем и данными о запчастях
type InventoryService struct {
	inventory_v1.UnimplementedInventoryServiceServer

	parts *PartStorage
}

//GetPart отправляет запрос на предоставление данных о запчасте
func (s *InventoryService) GetPart(_ context.Context, req *inventory_v1.GetPartRequest) (*inventory_v1.GetPartResponse, error) {
	part := s.parts.GetPart(req.Uuid)
	if part == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("part uuid %s is not found", req.Uuid))
	}

	return &inventory_v1.GetPartResponse{
		Part: part,
	}, nil
}

func (s *InventoryService) ListParts(_ context.Context, req *inventory_v1.ListPartsRequest) (*inventory_v1.ListPartsResponse, error) {
	filter := req.GetFilter()
	if filter == nil {
		// фильтр отсутствует — возвращаем все детали
		return &inventory_v1.ListPartsResponse{Parts: s.parts.GetAllParts()}, nil
	}
	uuids := filter.GetUuids()

	var parts []*inventory_v1.Part

	// 1. Создаём список запчастей по uuids.
	// Если uuids не передано возвращает список всех запчастей из хранилища
	if len(uuids) != 0 {
		parts = s.parts.GetParts(uuids)
	} else {
		parts = s.parts.GetAllParts()
	}

	// 2. Фильтр по именам
	if names := filter.GetNames(); len(names) > 0 {
		set := toSet(names)
		parts = filterParts(parts, func(p *inventory_v1.Part) bool {
			_, ok := set[p.GetName()]
			return ok
		})
	}

	// 3. Фильтр по категориям
	if categories := filter.GetCategories(); len(categories) > 0 {
		set := toSet(categories)
		parts = filterParts(parts, func(p *inventory_v1.Part) bool {
			_, ok := set[p.GetCategory()]
			return ok
		})
	}

	// 4. Фильтр по странам производителя
	if countries := filter.GetManufacturerCountries(); len(countries) > 0 {
		set := toSet(countries)
		parts = filterParts(parts, func(p *inventory_v1.Part) bool {
			if m := p.GetManufacturer(); m != nil {
				_, ok := set[m.GetCountry()]
				return ok
			}
			return false
		})
	}

	// 5. Фильтр по тегам
	if tags := filter.GetTags(); len(tags) > 0 {
		set := toSet(tags)
		parts = filterParts(parts, func(p *inventory_v1.Part) bool {
			for _, tag := range p.GetTags() {
				if _, ok := set[tag]; ok {
					return true
				}
			}
			return false
		})
	}

	return &inventory_v1.ListPartsResponse{Parts: parts}, nil
}

// toSet преобразует слайс элементов в множество (map[T]struct{}).
// Требуется comparable, чтобы ключи можно было сравнивать.
func toSet[T comparable](items []T) map[T]struct{} {
	set := make(map[T]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}

// filterParts создаёт новый слайс из элементов, удовлетворяющих предикату.
func filterParts(parts []*inventory_v1.Part, keep func(*inventory_v1.Part) bool) []*inventory_v1.Part {
	res := make([]*inventory_v1.Part, 0, len(parts))
	for _, p := range parts {
		if keep(p) {
			res = append(res, p)
		}
	}
	return res
}

func panicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC in %s: %v\n%s", info.FullMethod, r, string(debug.Stack()))
				err = status.Error(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		fmt.Printf("failed to listen: %v\n", err)
		return
	}

	defer func() {
		err := lis.Close()
		if err != nil {
			fmt.Printf("failed to close listener: %v\n", err)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer(
		grpc.UnaryInterceptor(panicRecoveryInterceptor()),
	)

	//Регистрируем Inventory сервис
	service := &InventoryService{
		parts: NewPartStorage(),
	}

	// Наполняем тестовыми данными
	service.parts.UpdatePart(&inventory_v1.Part{
		Uuid:          "11111111-1111-1111-1111-111111111111",
		Name:          "Сопло маршевое",
		Price:         1500.0,
		StockQuantity: 5,
		Category:      inventory_v1.Category_CATEGORY_ENGINE,
		Tags:          []string{"engine", "main"},
	})
	service.parts.UpdatePart(&inventory_v1.Part{
		Uuid:          "22222222-2222-2222-2222-222222222222",
		Name:          "Иллюминатор стандартный",
		Price:         300.0,
		StockQuantity: 12,
		Category:      inventory_v1.Category_CATEGORY_PORTHOLE,
		Tags:          []string{"porthole", "window"},
	})

	inventory_v1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		fmt.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err := s.Serve(lis)
		if err != nil {
			fmt.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	//Gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
