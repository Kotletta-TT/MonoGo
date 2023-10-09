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
	d := &Database{db: db}
	logger.Info("Create table")
	d.initTable()
	return d, nil
}

func (d *Database) initTable() {
	if _, err := d.db.Exec(`CREATE TABLE IF NOT EXISTS metrics (
		name VARCHAR(255) NOT NULL,
		mtype VARCHAR(10) NOT NULL,
		value DOUBLE PRECISION,
		delta INTEGER,
		PRIMARY KEY (name)
	);`); err != nil {
		panic(err)
	}
}

func (d *Database) StoreGaugeMetric(name string, value float64) {
	tx, err := d.db.Begin()
	if err != nil {
		panic(err)
	}
	if _, err := tx.Exec(`INSERT INTO metrics (
		name, 
		mtype, 
		value) VALUES ($1, $2, $3)
		ON CONFLICT (name) 
		DO UPDATE SET value = $3`,
		name,
		"gauge",
		value); err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func (d *Database) StoreCounterMetric(name string, value int64) {
	tx, err := d.db.Begin()
	if err != nil {
		panic(err)
	}
	if _, err := tx.Exec(`INSERT INTO metrics (
		name, 
		mtype, 
		delta) VALUES ($1, $2, $3)
		ON CONFLICT (name) 
		DO UPDATE SET delta = $3`,
		name,
		"counter",
		value); err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func (d *Database) GetGaugeMetric(name string) (float64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		panic(err)
	}
	var value float64
	err = tx.QueryRow(`SELECT value FROM metrics WHERE name = $1`, name).Scan(&value)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	return value, nil
}

func (d *Database) GetCounterMetric(name string) (int64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		panic(err)
	}
	var value int64
	err = tx.QueryRow(`SELECT delta FROM metrics WHERE name = $1`, name).Scan(&value)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	return value, nil
}

func (d *Database) GetAllMetrics() map[string]*shared.Metrics {
	tx, err := d.db.Begin()
	if err != nil {
		panic(err)
	}
	rows, err := tx.Query(`SELECT name, mtype, value, delta FROM metrics`)
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	defer rows.Close()
	metrics := make(map[string]*shared.Metrics)
	for rows.Next() {
		var name, mtype string
		var value sql.NullFloat64
		var delta sql.NullInt64
		err = rows.Scan(&name, &mtype, &value, &delta)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
		metrics[name] = &shared.Metrics{
			ID:    name,
			MType: mtype,
		}
		if value.Valid {
			metrics[name].Value = &value.Float64
		}
		if delta.Valid {
			metrics[name].Delta = &delta.Int64
		}
	}
	if err = tx.Commit(); err != nil {
		panic(err)
	}
	return metrics
}

func (d *Database) LoadFromFile() (map[string]*shared.Metrics, error) {
	return nil, nil
}

func (d *Database) Stash() {}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) HealthCheck(ctx context.Context) error {
	return d.db.PingContext(ctx)
}
