package inventory

import repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"

func (r *repository) Init() {
	// Наполняем тестовыми данными
	r.UpdatePart(&repoModel.Part{
		UUID:          "11111111-1111-1111-1111-111111111111",
		Name:          "Сопло маршевое",
		Price:         1500.0,
		StockQuantity: 5,
		Category:      repoModel.Category_CATEGORY_ENGINE,
		Manufacturer: &repoModel.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		},
		Tags: []string{"engine", "main"},
	})
	r.UpdatePart(&repoModel.Part{
		UUID:          "22222222-2222-2222-2222-222222222222",
		Name:          "Иллюминатор стандартный",
		Price:         300.0,
		StockQuantity: 12,
		Category:      repoModel.Category_CATEGORY_PORTHOLE,
		Manufacturer: &repoModel.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		},
		Tags: []string{"porthole", "window"},
	})
	r.UpdatePart(&repoModel.Part{
		UUID:          "33333333-3333-3333-3333-333333333333",
		Name:          "Иллюминатор квадратный",
		Price:         600.0,
		StockQuantity: 2,
		Category:      repoModel.Category_CATEGORY_PORTHOLE,
		Manufacturer:  nil,
		Tags:          nil,
	})
}

func (r *repository) UpdatePart(part *repoModel.Part) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[part.UUID] = *part
}
