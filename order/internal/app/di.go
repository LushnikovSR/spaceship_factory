package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiOrderV1 "github.com/LushnikovSR/spaceship_factory/order/internal/api/order/v1"
	grpcClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc"
	inventoryClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/LushnikovSR/spaceship_factory/order/internal/client/grpc/payment/v1"
	config "github.com/LushnikovSR/spaceship_factory/order/internal/config"
	repository "github.com/LushnikovSR/spaceship_factory/order/internal/repository"
	repositoryOrder "github.com/LushnikovSR/spaceship_factory/order/internal/repository/order"
	service "github.com/LushnikovSR/spaceship_factory/order/internal/service"
	serviceOrder "github.com/LushnikovSR/spaceship_factory/order/internal/service/order"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/closer"
	"github.com/LushnikovSR/spaceship_factory/platform/pkg/logger"
	pgmigrator "github.com/LushnikovSR/spaceship_factory/platform/pkg/migrator/pg"
	orderV1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/openapi/order/v1"
	inventory_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/LushnikovSR/spaceship_factory/shared/pkg/proto/payment/v1"
)

const (
	orderServerPrefixPath = "/api/v1"
)

type diContainer struct {
	orderV1API          orderV1.Handler
	orderService        service.OrderService
	orderResository     repository.OrderRepository
	pgxpoolClient       *pgxpool.Pool
	inventoryGRPCClient grpcClient.InventoryClient
	paymentGRPCClient   grpcClient.PaymentClient
	chiRouter           chi.Router
	orderServer         http.Handler
	pgMigrator          *pgmigrator.Migrator
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = apiOrderV1.NewAPI(d.OrderService(ctx))
	}
	return d.orderV1API
}

func (d *diContainer) InventoryClient(ctx context.Context) grpcClient.InventoryClient {
	if d.inventoryGRPCClient == nil {
		conn, err := grpc.NewClient(config.AppConfig().InventoryGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error(ctx, "failed to create new Inventory grpcClient", zap.Error(err))
			return nil
		}
		closer.AddNamed("InventoryGRPCClient", func(ctx context.Context) error {
			return conn.Close()
		})

		clientConn := inventory_v1.NewInventoryServiceClient(conn)
		d.inventoryGRPCClient = inventoryClient.NewClient(clientConn)
	}
	return d.inventoryGRPCClient
}

func (d *diContainer) PaymentClient(ctx context.Context) grpcClient.PaymentClient {
	if d.paymentGRPCClient == nil {
		conn, err := grpc.NewClient(config.AppConfig().PaymentGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			logger.Error(ctx, "failed to create new Payment grpcClient", zap.Error(err))
			return nil
		}
		closer.AddNamed("PaymentGRPCClient", func(ctx context.Context) error {
			return conn.Close()
		})

		clientConn := payment_v1.NewPaymentServiceClient(conn)
		d.paymentGRPCClient = paymentClient.NewClient(clientConn)
	}
	return d.paymentGRPCClient
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = serviceOrder.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(ctx),
			d.PaymentClient(ctx),
		)
	}
	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	var err error
	if d.orderResository == nil {
		d.orderResository, err = repositoryOrder.NewRepository(d.PgxpoolClient(ctx))
		if err != nil {
			logger.Error(ctx, "failed to create Order repository", zap.Error(err))
			return nil
		}
	}
	return d.orderResository
}

func (d *diContainer) PgxpoolClient(ctx context.Context) *pgxpool.Pool {
	if d.pgxpoolClient == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			logger.Error(ctx, "failed to connect to postgres", zap.Error(err))
			return nil
		}
		closer.AddNamed("pgxpool Pool", func(ctx context.Context) error {
			pool.Close()
			return nil
		})
		d.pgxpoolClient = pool
	}

	d.PGMigrator()

	return d.pgxpoolClient
}

func (d *diContainer) ChiRouter(ctx context.Context) chi.Router {
	if d.chiRouter == nil {
		r := chi.NewRouter()
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Timeout(10 * time.Second))
		r.Use(render.SetContentType(render.ContentTypeJSON))

		orderServer := d.OrderServer(ctx)
		r.Mount("/", orderServer)

		d.chiRouter = r
	}
	return d.chiRouter
}

func (d *diContainer) OrderServer(ctx context.Context) http.Handler {
	if d.orderServer == nil {
		apiHandler := d.OrderV1API(ctx)
		orderServer, err := orderV1.NewServer(apiHandler, orderV1.WithPathPrefix(orderServerPrefixPath))
		if err != nil {
			logger.Error(ctx, "failed to create OpenAPI server", zap.Error(err))
			return nil
		}

		d.orderServer = orderServer
	}
	return d.orderServer
}

func (d *diContainer) PGMigrator() {
	db := stdlib.OpenDB(*d.pgxpoolClient.Config().ConnConfig.Copy())
	d.pgMigrator = pgmigrator.Init(db, config.AppConfig().Postgres.MigrationDir())

	closer.AddNamed("pgMigrator", func(ctx context.Context) error {
		return db.Close()
	})
}
