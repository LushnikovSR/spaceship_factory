package integration

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
)

var _ = Describe("InventorySerivce", func() {
	var (
		ctx             context.Context
		cancel          context.CancelFunc
		inventoryClient inventoryV1.InventoryServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(suiteCtx)

		// Создаём gRPC клиент
		conn, err := grpc.NewClient(
			env.App.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		Expect(err).ToNot(HaveOccurred(), "ожидали успешное подключение к gRPC приложению")

		inventoryClient = inventoryV1.NewInventoryServiceClient(conn)

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

			req := inventoryV1.GetPartRequest{
				Uuid: uuid,
			}

			resp, err := inventoryClient.GetPart(ctx, &req)
			Expect(err).NotTo(HaveOccurred(), "ожидали успешное получение детали из коллекции parts")
			Expect(resp.Part.Uuid).To(Equal(uuid))
		})
	})
})
