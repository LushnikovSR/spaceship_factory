package inventory

import (
	"context"
	"fmt"
	"log/slog"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/converter"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartsFilter) (parts []*model.Part, err error) {
	if filter == nil {
		// фильтр отсутствует — возвращаем все детали
		repoParts, err := r.GetAllParts(ctx)
		if err != nil {
			return nil, err
		}
		parts = converter.ToModelParts(repoParts)
		return parts, nil
	}
	uuids := filter.Uuids

	var repoParts []repoModel.Part

	// 1. Создаётся список запчастей по uuids.
	// Если uuids это пустой список, возвращаются все запчасти
	repoParts, err = r.GetParts(ctx, uuids)
	if err != nil {
		return nil, err
	}

	// 2. Фильтр по именам
	if names := filter.Names; len(names) > 0 {
		set := toSet(names)
		repoParts = filterParts(repoParts, func(p repoModel.Part) bool {
			_, ok := set[p.Name]
			return ok
		})
	}

	// 3. Фильтр по категориям
	if categories := filter.Categories; len(categories) > 0 {
		set := toSet(categories)
		repoParts = filterParts(repoParts, func(p repoModel.Part) bool {
			_, ok := set[model.Category(p.Category)]
			return ok
		})
	}

	// 4. Фильтр по странам производителя
	if countries := filter.ManufacturerCountries; len(countries) > 0 {
		set := toSet(countries)
		repoParts = filterParts(repoParts, func(p repoModel.Part) bool {
			if m := p.Manufacturer; m != nil {
				_, ok := set[m.Country]
				return ok
			}
			return false
		})
	}

	// 5. Фильтр по тегам
	if tags := filter.Tags; len(tags) > 0 {
		set := toSet(tags)
		repoParts = filterParts(repoParts, func(p repoModel.Part) bool {
			for _, tag := range p.Tags {
				if _, ok := set[tag]; ok {
					return true
				}
			}
			return false
		})
	}

	parts = converter.ToModelParts(repoParts)

	return parts, nil
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
func filterParts(parts []repoModel.Part, keep func(repoModel.Part) bool) []repoModel.Part {
	res := make([]repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if keep(p) {
			res = append(res, p)
		}
	}
	return res
}

func (r *repository) GetParts(ctx context.Context, uuids []string) ([]repoModel.Part, error) {
	// Ограничение на количество идентификаторов, чтобы не замедлять запрос и не выйти за размер документа MongoDB
	const maxIDs = 500
	if len(uuids) > maxIDs {
		return nil, fmt.Errorf("requested %d IDs, max allowed %d", len(uuids), maxIDs)
	}

	// Исключение дубликатов из запроса к базе
	seen := make(map[string]struct{}, len(uuids))
	uniqueIDs := make([]string, 0)
	for _, id := range uuids {
		if _, ok := seen[id]; !ok && id != "" {
			uniqueIDs = append(uniqueIDs, id)
		}
	}
	uuids = uniqueIDs

	// Возвращаются все детали при len(uuids) == 0
	if len(uuids) == 0 {
		parts, err := r.GetAllParts(ctx)
		if err != nil {
			return nil, err
		}
		return parts, nil
	}

	// Получение деталей из базы данных по uuids
	filter := bson.M{"_id": bson.M{"$in": uuids}}
	cursor, err := r.data.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find parts: %w", err)
	}

	defer func() {
		cerr := cursor.Close(context.Background())
		if cerr != nil {
			slog.Warn("failed to close cursor", "error", cerr)
		}
	}()

	parts := make([]repoModel.Part, 0, len(uuids))
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	return parts, nil
}

// GetAll получает все запчасти хранящиеся в базе
func (r *repository) GetAllParts(ctx context.Context) ([]repoModel.Part, error) {
	cursor, err := r.data.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			slog.Info("failed to close cursor", "error", cerr)
		}
	}()

	var parts []repoModel.Part
	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}

	return parts, nil
}
