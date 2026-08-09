package pgmigrator

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/LushnikovSR/spaceship_factory/platform/pkg/logger"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

type Migrator struct {
	db  *sql.DB
	dir string
}

func Init(db *sql.DB, dir string) *Migrator {
	return &Migrator{
		db:  db,
		dir: dir,
	}
}

func (m *Migrator) Run(ctx context.Context) error {
	start := time.Now()

	versionBefore, err := goose.GetDBVersion(m.db)
	if err != nil {
		logger.Warn(ctx, "не удалось получить версию БД до миграции", zap.Error(err))
	}

	err = goose.UpContext(ctx, m.db, m.dir)
	if err != nil {
		logger.Error(ctx, "миграции не выполнены",
			zap.Int64("version_before", versionBefore),
			zap.Error(err),
			zap.Duration("elapsed", time.Since(start)),
		)
		return fmt.Errorf("ошибка миграций: %w", err)
	}

	versionAfter, err := goose.GetDBVersion(m.db)
	if err != nil {
		logger.Warn(ctx, "не удалось получить версию БД после миграции", zap.Error(err))
	}

	logger.Info(ctx, "миграции успешно применены",
		zap.Int64("version_before", versionBefore),
		zap.Int64("version_after", versionAfter),
		zap.Duration("duration", time.Since(start)),
	)
	return nil
}

func (m *Migrator) Up(ctx context.Context) error {
	return goose.UpContext(ctx, m.db, m.dir)
}

func (m *Migrator) Down(ctx context.Context) error {
	return goose.DownContext(ctx, m.db, m.dir)
}

func (m *Migrator) Status(ctx context.Context) error {
	return goose.StatusContext(ctx, m.db, m.dir)
}
