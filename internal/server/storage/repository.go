package storage

import (
	"context"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage/memory"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage/postgres"
	"github.com/Kotletta-TT/MonoGo/internal/shared"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository interface {
	StoreGaugeMetric(name string, value float64)
	StoreCounterMetric(name string, value int64)
	GetGaugeMetric(name string) (float64, error)
	GetCounterMetric(name string) (int64, error)
	GetAllMetrics() map[string]*shared.Metrics
	LoadFromFile() (map[string]*shared.Metrics, error)
	Stash()
	Close()
	HealthCheck(ctx context.Context) error
}

func GetRepo(cfg *config.Config) Repository {
	repo, err := postgres.New(cfg)
	if cfg.DatabaseDSN == "" || err != nil {
		logger.Infof("Connect to database error: %s", err)
		logger.Info("Repo: MemoryStorage")
		return memory.New(cfg)
	}
	logger.Info("Repo: PostgreSQL")
	return repo
}
