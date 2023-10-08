package postgres

import (
	"context"
	"database/sql"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/shared"
)

type Database struct {
	db *sql.DB
}

func New(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}

func (d *Database) StoreGaugeMetric(name string, value float64) {
	tx, err := d.db.Begin()
	if err != nil {
		logger.Fatal(err)
	}
	
	tx.Exec()
}

func (d *Database) StoreCounterMetric(name string, value int64) {
	panic("not implemented") // TODO: Implement
}

func (d *Database) GetGaugeMetric(name string) (float64, error) {
	panic("not implemented") // TODO: Implement
}

func (d *Database) GetCounterMetric(name string) (int64, error) {
	panic("not implemented") // TODO: Implement
}

func (d *Database) GetAllMetrics() map[string]*shared.Metrics {
	panic("not implemented") // TODO: Implement
}

func (d *Database) LoadFromFile() (map[string]*shared.Metrics, error) {
	panic("not implemented") // TODO: Implement
}

func (d *Database) Stash() {
	panic("not implemented") // TODO: Implement
}

func (d *Database) Close() {
	panic("not implemented") // TODO: Implement
}

func (d *Database) HealthCheck(ctx context.Context) error {
	panic("not implemented") // TODO: Implement
}
