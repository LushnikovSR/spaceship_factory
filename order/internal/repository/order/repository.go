package order

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	migrator "github.com/LushnikovSR/spaceship_factory/order/internal/migrator"
	def "github.com/LushnikovSR/spaceship_factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	data *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) (*repository, error) {
	db := stdlib.OpenDB(*pool.Config().ConnConfig.Copy()) // Создаётся копия *pgxpool.Pool и приводится к типу *sql.DB
	defer func() {
		cerr := db.Close()
		if cerr != nil {
			slog.Warn("failed to close cursor", "error", cerr)
		}
	}()

	// Инициализируем мигратор
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	migratorRunner := migrator.NewMigrator(db, migrationsDir)

	err := migratorRunner.Up()
	if err != nil {
		return nil, fmt.Errorf("database migration error: %w", err)
	}

	return &repository{
		data: pool,
	}, nil
}

// ConnectPostgres создаёт соединение с базой данных. Требует вызова defer *sql.DB.Close()
func ConnectPostgres(ctx context.Context) (*pgxpool.Pool, error) {
	// Закружаем переменные окружения
	err := godotenv.Load("deploy/compose/order/.env")
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	dbURI := os.Getenv("POSTGRES_URI")

	// Создаём соединение с базой данных
	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	// Проверяем что соединение с базой данных установлено
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("database is unavailable: %w", err)
	}

	return pool, nil
}
