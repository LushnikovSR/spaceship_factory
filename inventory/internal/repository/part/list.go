package inventory

import (
	"context"

	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
	converter "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/converter"
	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	if filter == nil {
		// фильтр отсутствует — возвращаем все детали
		repoParts := r.GetAllParts()
		parts := toModelParts(repoParts)
		return parts, nil
	}
	uuids := filter.Uuids

	var repoParts []*repoModel.Part

	// 1. Создаём список запчастей по uuids.
	// Если uuids не передано возвращает список всех запчастей из хранилища
	if len(uuids) != 0 {
		repoParts = r.GetParts(uuids)
	} else {
		repoParts = r.GetAllParts()
	}

	// 2. Фильтр по именам
	if names := filter.Names; len(names) > 0 {
		set := toSet(names)
		repoParts = filterParts(repoParts, func(p *repoModel.Part) bool {
			_, ok := set[p.Name]
			return ok
		})
	}

	// 3. Фильтр по категориям
	if categories := filter.Categories; len(categories) > 0 {
		set := toSet(categories)
		repoParts = filterParts(repoParts, func(p *repoModel.Part) bool {
			_, ok := set[model.Category(p.Category)]
			return ok
		})
	}

	// 4. Фильтр по странам производителя
	if countries := filter.ManufacturerCountries; len(countries) > 0 {
		set := toSet(countries)
		repoParts = filterParts(repoParts, func(p *repoModel.Part) bool {
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
		repoParts = filterParts(repoParts, func(p *repoModel.Part) bool {
			for _, tag := range p.Tags {
				if _, ok := set[tag]; ok {
					return true
				}
			}
			return false
		})
	}

	parts := toModelParts(repoParts)

	return parts, nil
}

func toModelParts(repoParts []*repoModel.Part) []*model.Part {
	parts := make([]*model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		part := converter.RepoModelToModel(repoPart)
		parts = append(parts, part)
	}
	return parts
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
func filterParts(parts []*repoModel.Part, keep func(*repoModel.Part) bool) []*repoModel.Part {
	res := make([]*repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if keep(p) {
			res = append(res, p)
		}
	}
	return res
}

func (r *repository) GetParts(uuids []string) []*repoModel.Part {
	seen := make(map[string]struct{}, len(uuids))
	parts := make([]*repoModel.Part, 0, len(uuids))

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, id := range uuids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		if part, ok := r.data[id]; ok {
			parts = append(parts, &part)
		}
	}

	return parts
}

// GetAllParts возвращает информацию о всех запчастях в хранилище.
func (r *repository) GetAllParts() []*repoModel.Part {
	parts := make([]*repoModel.Part, 0)

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, part := range r.data {
		parts = append(parts, &part)
	}

	return parts
}
