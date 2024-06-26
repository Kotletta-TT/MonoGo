// Package storage provides repository for server
package storage

import (
	"context"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage/memory"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Repository Данный интерфейс отображает на мой взгляд самый быстрый способ получения и заполнения данных,
// Сокращает объем кода, и использует мало памяти, засчет эффективного переиспользования структур и типов данных.
type Repository interface {
	StoreMetric(metric *common.Metrics) error
	StoreBatchMetric(metrics []*common.Metrics) ([]*common.Metrics, error)
	GetMetric(metric *common.Metrics) error
	GetListMetrics() ([]*common.Metrics, error)
	HealthCheck(ctx context.Context) error
	Close()
}

// GetRepo returns a Repository based on the provided Config.
//
// It takes a pointer to a Config as its parameter and returns a Repository.
func GetRepo(ctx context.Context, cfg *config.Config) Repository {
	repo, err := postgres.New(ctx, cfg)
	if cfg.DatabaseDSN == "" || err != nil {
		logger.Infof("Connect to database error: %s", err)
		logger.Info("Repo: MemoryStorage")
		repo := memory.New(cfg)
		go repo.Stash()
		return repo
	}
	logger.Info("Repo: PostgreSQL")
	return repo
}
