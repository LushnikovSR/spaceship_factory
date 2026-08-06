//go:build integration

package integration

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	repoModel "github.com/LushnikovSR/spaceship_factory/inventory/internal/repository/model"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventorySerivce", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventory_v1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventory_v1.NewInventoryServiceClient(conn)

	})

	AfterEach(func() {
		// Чистим коллекцию после теста
		err := env.ClearPartsCollection(ctx)
		Expect(err).NotTo(HaveOccurred(), "ожидали успешную очистку коллекции parts")

		cancel()
	})

	Describe("GetPart", func() {
		It("must successfully returns information about a part by its uuid", func() {
			uuid, err := env.InsertTestPart(ctx)
			Expect(err).NotTo(HaveOccurred(), "ожидали успешное добавление детали в коллекцию parts")

			req := inventory_v1.GetPartRequest{
				Uuid: uuid,
			}

			resp, err := inventoryClient.GetPart(ctx, &req)
			Expect(err).NotTo(HaveOccurred(), "ожидали успешное получение детали из коллекции parts")
			Expect(resp.Part.Uuid).To(Equal(uuid))
		})
	})

	Describe("ListParts", func() {
		var (
			amountCategories = 4
			n                = 10
			uuids            = make([]string, 0, n)
			names            = make([]string, 0, n)
			categories       = make([]inventory_v1.Category, 0, n)
			mCountries       = make([]string, 0, n)
			tags             = make([]string, 0, n)

			filter *inventory_v1.PartsFilter
		)

		BeforeEach(func() {
			//Вставляем искомую тестовую деталь и задаём параметры поиска
			uuid, err := env.InsertTestPart(ctx)
			Expect(err).NotTo(HaveOccurred(), "expected that the part would be successfully added to the 'parts' collection")

			uuids = append(uuids, uuid)

			part := env.GetTestPart()

			names = append(names, part.Name)
			categories = append(categories, inventory_v1.Category(part.Category))
			mCountries = append(mCountries, part.Manufacturer.Country)
			tags = append(tags, part.Tags...)

			filter = &inventory_v1.PartsFilter{
				Uuids:                 uuids,
				Names:                 names,
				Categories:            categories,
				ManufacturerCountries: mCountries,
				Tags:                  tags,
			}

			// Вставляем тестовые детали
			for i := 1; i < n; i++ {
				numStr := strconv.Itoa(i)
				p := repoModel.Part{
					Name:          testPart.Name + "_" + numStr,
					Price:         testPart.Price + float64(i*n*n),
					StockQuantity: testPart.StockQuantity + int64(i),
					Category:      repoModel.Category(rand.Intn(amountCategories)),
				}

				_, err := env.InsertTestPartWithData(ctx, p)
				Expect(err).NotTo(HaveOccurred(), "expected that the collection of 'parts' would be successfully populated with data")
			}
		})

		It("must successfully returns information about a part by filter", func() {
			req := &inventory_v1.ListPartsRequest{
				Filter: filter,
			}

			resp, err := inventoryClient.ListParts(ctx, req)
			Expect(err).NotTo(HaveOccurred(), "expected to receive the part from the 'parts' collection")
			Expect(resp.Parts).ShouldNot(BeEmpty(), fmt.Sprintf("expected to receive the part from the 'parts' collection with uuid: %s", uuids[0]))
			Expect(resp.Parts[0].Uuid).To(Equal(uuids[0]))
		})
	})
})
