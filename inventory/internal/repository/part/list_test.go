package inventory

import (
	model "github.com/LushnikovSR/spaceship_factory/inventory/internal/model"
)

func (s *RepositorySuite) TestListParts_FullFilter() {
	var (
		partUuids                 = []string{"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222"}
		partNmaes                 = []string{"Сопло маршевое", "Иллюминатор стандартный"}
		partCategories            = []model.Category{model.Category_CATEGORY_ENGINE, model.Category_CATEGORY_PORTHOLE}
		partManufacturerCountries = []string{"Germany"}
		partTags                  = []string{"engine", "window"}

		filter = &model.PartsFilter{
			Uuids:                 partUuids,
			Names:                 partNmaes,
			Categories:            partCategories,
			ManufacturerCountries: partManufacturerCountries,
			Tags:                  partTags,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
			{
				UUID:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  manufacturer,
				Tags:          []string{"porthole", "window"},
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_NillFilter() {
	var (
		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
			{
				UUID:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  manufacturer,
				Tags:          []string{"porthole", "window"},
			},
			{
				UUID:          "33333333-3333-3333-3333-333333333333",
				Name:          "Иллюминатор квадратный",
				Price:         600.0,
				StockQuantity: 2,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  nil,
				Tags:          nil,
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, nil)
	s.Require().NoError(err)
	s.Require().ElementsMatch(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_WithNonExistentPartUuidAmongPartUuids() {
	var (
		nonExistentPartUUID       = "non-existent"
		partUuids                 = []string{"11111111-1111-1111-1111-111111111111", nonExistentPartUUID}
		partNmaes                 = []string{"Сопло маршевое", "Иллюминатор стандартный"}
		partCategories            = []model.Category{model.Category_CATEGORY_ENGINE, model.Category_CATEGORY_PORTHOLE}
		partManufacturerCountries = []string{"Germany"}
		partTags                  = []string{"engine", "window"}

		filter = &model.PartsFilter{
			Uuids:                 partUuids,
			Names:                 partNmaes,
			Categories:            partCategories,
			ManufacturerCountries: partManufacturerCountries,
			Tags:                  partTags,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_WithNonExistentPartUUID() {
	var (
		nonExistentPartUUID       = "non-existent"
		partUuids                 = []string{nonExistentPartUUID}
		partNmaes                 = []string{"Сопло маршевое", "Иллюминатор стандартный"}
		partCategories            = []model.Category{model.Category_CATEGORY_ENGINE, model.Category_CATEGORY_PORTHOLE}
		partManufacturerCountries = []string{"Germany"}
		partTags                  = []string{"engine", "window"}

		filter = &model.PartsFilter{
			Uuids:                 partUuids,
			Names:                 partNmaes,
			Categories:            partCategories,
			ManufacturerCountries: partManufacturerCountries,
			Tags:                  partTags,
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Empty(parts)
}

func (s *RepositorySuite) TestListParts_WithNonExistentPartName() {
	var (
		nonExistentPartName       = "non-existent"
		partUuids                 = []string{"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222"}
		partNmaes                 = []string{nonExistentPartName}
		partCategories            = []model.Category{model.Category_CATEGORY_ENGINE, model.Category_CATEGORY_PORTHOLE}
		partManufacturerCountries = []string{"Germany"}
		partTags                  = []string{"engine", "window"}

		filter = &model.PartsFilter{
			Uuids:                 partUuids,
			Names:                 partNmaes,
			Categories:            partCategories,
			ManufacturerCountries: partManufacturerCountries,
			Tags:                  partTags,
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Empty(parts)
}

func (s *RepositorySuite) TestListParts_WithOnlyUuidsInFilter() {
	var (
		partUuids = []string{"11111111-1111-1111-1111-111111111111", "22222222-2222-2222-2222-222222222222"}

		filter = &model.PartsFilter{
			Uuids: partUuids,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
			{
				UUID:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  manufacturer,
				Tags:          []string{"porthole", "window"},
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_WithDoubleUuidInFilter() {
	var (
		partUuids = []string{"11111111-1111-1111-1111-111111111111", "11111111-1111-1111-1111-111111111111"}

		filter = &model.PartsFilter{
			Uuids: partUuids,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_WithOnlyNamesInFilter() {
	var (
		partNmaes = []string{"Сопло маршевое", "Иллюминатор стандартный"}

		filter = &model.PartsFilter{
			Names: partNmaes,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
			{
				UUID:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  manufacturer,
				Tags:          []string{"porthole", "window"},
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().ElementsMatch(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_WithOnlyCategoriesInFilter() {
	var (
		partCategories = []model.Category{model.Category_CATEGORY_ENGINE, model.Category_CATEGORY_PORTHOLE}

		filter = &model.PartsFilter{
			Categories: partCategories,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
			{
				UUID:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  manufacturer,
				Tags:          []string{"porthole", "window"},
			},
			{
				UUID:          "33333333-3333-3333-3333-333333333333",
				Name:          "Иллюминатор квадратный",
				Price:         600.0,
				StockQuantity: 2,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  nil,
				Tags:          nil,
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().ElementsMatch(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_WithOnlyManufacturerCountriesInFilter() {
	var (
		partManufacturerCountries = []string{"Germany"}

		filter = &model.PartsFilter{
			ManufacturerCountries: partManufacturerCountries,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
			{
				UUID:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  manufacturer,
				Tags:          []string{"porthole", "window"},
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedParts, parts)
}

func (s *RepositorySuite) TestListParts_WithOnlyTagsInFilter() {
	var (
		partTags = []string{"engine", "window"}

		filter = &model.PartsFilter{
			Tags: partTags,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Biscuit",
			Country: "Germany",
			Website: "financialharness.info",
		}

		expectedParts = []*model.Part{
			{
				UUID:          "11111111-1111-1111-1111-111111111111",
				Name:          "Сопло маршевое",
				Price:         1500.0,
				StockQuantity: 5,
				Category:      model.Category_CATEGORY_ENGINE,
				Manufacturer:  manufacturer,
				Tags:          []string{"engine", "main"},
			},
			{
				UUID:          "22222222-2222-2222-2222-222222222222",
				Name:          "Иллюминатор стандартный",
				Price:         300.0,
				StockQuantity: 12,
				Category:      model.Category_CATEGORY_PORTHOLE,
				Manufacturer:  manufacturer,
				Tags:          []string{"porthole", "window"},
			},
		}
	)
	s.repository.Init()

	parts, err := s.repository.ListParts(s.ctx, filter)
	println(parts)
	s.Require().NoError(err)
	s.Require().ElementsMatch(expectedParts, parts)
}
