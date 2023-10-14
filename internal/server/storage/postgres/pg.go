package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/shared"
)

type Database struct {
	db    *sql.DB
	cache map[string]*shared.Metrics
	mu    sync.Mutex
}

const GAUGE = "gauge"
const COUNTER = "counter"

const initCacheQuery = `SELECT name FROM metrics`

const createTableQuery = `
	CREATE TABLE 
	IF NOT EXISTS metrics (
	name VARCHAR(255) NOT NULL,
	mtype VARCHAR(10) NOT NULL,
	value DOUBLE PRECISION,
	delta BIGINT,
	PRIMARY KEY (name));`

const insertQueryGauge = `
	INSERT INTO metrics (name, mtype, value)
	VALUES ($1, $2, $3)
	ON CONFLICT (name)
	DO UPDATE SET value = EXCLUDED.value`

const insertQueryCounter = `
	INSERT INTO metrics (name, mtype, delta) 
	VALUES ($1, $2, $3)
	ON CONFLICT (name)
	DO UPDATE SET delta = EXCLUDED.delta`

const insertBatchQuery = `
	INSERT INTO metrics (name, mtype, value, delta)
	VALUES %s
	ON CONFLICT (name)
	DO UPDATE SET (value, delta) = (EXCLUDED.value, EXCLUDED.delta)
`

const selectQueryGauge = `SELECT value FROM metrics WHERE name = $1 AND mtype = $2`
const selectQueryCounter = `SELECT delta FROM metrics WHERE name = $1 AND mtype = $2`
const selectAllMetrics = `SELECT name, mtype, value, delta FROM metrics`

func New(cfg *config.Config) (*Database, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	d := &Database{db: db, cache: make(map[string]*shared.Metrics), mu: sync.Mutex{}}
	logger.Info("Create table")
	err = d.initTable()
	if err != nil {
		return nil, err
	}
	err = d.db.Ping()
	if err != nil {
		return nil, err
	}
	err = d.initCache()
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (d *Database) initCache() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	cache, err := d.GetAllMetrics()
	if err != nil {
		return err
	}
	d.cache = cache
	return nil
}

func (d *Database) initTable() error {
	if _, err := d.db.Exec(createTableQuery); err != nil {
		return err
	}
	return nil
}

func (d *Database) StoreGaugeMetric(name string, value float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, err := d.db.Exec(insertQueryGauge, name, GAUGE, value)
	if err != nil {
		return err
	}
	d.cache[name] = &shared.Metrics{ID: name, MType: COUNTER, Value: &value}
	return nil
}

func (d *Database) StoreCounterMetric(name string, value int64) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	m, ok := d.cache[name]
	var newValue int64
	if ok {
		newValue = *m.Delta + value
	}
	_, err := d.db.Exec(insertQueryCounter, name, COUNTER, newValue)
	if err != nil {
		return err
	}
	d.cache[name] = &shared.Metrics{ID: name, MType: COUNTER, Delta: &newValue}
	return nil
}

func (d *Database) GetGaugeMetric(name string) (float64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	m, ok := d.cache[name]
	if ok {
		return *m.Value, nil
	}
	row := d.db.QueryRow(selectQueryGauge, name, GAUGE)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var value float64
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (d *Database) GetCounterMetric(name string) (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	m, ok := d.cache[name]
	if ok {
		return *m.Delta, nil
	}
	row := d.db.QueryRow(selectQueryCounter, name, COUNTER)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var delta int64
	err := row.Scan(&delta)
	if err != nil {
		return 0, err
	}
	return delta, nil
}

func (d *Database) GetAllMetrics() (map[string]*shared.Metrics, error) {
	rows, err := d.db.Query(selectAllMetrics)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	metrics := make(map[string]*shared.Metrics)
	for rows.Next() {
		var name, mtype string
		var value sql.NullFloat64
		var delta sql.NullInt64
		err = rows.Scan(&name, &mtype, &value, &delta)
		if err != nil {
			return nil, err
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
	return metrics, nil
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

func (d *Database) StoreBatchMetric(metricSlice []*shared.Metrics) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for _, metric := range metricSlice {
		cacheMetric, ok := d.cache[metric.ID]
		if !ok {
			d.cache[metric.ID] = metric
			continue
		}
		if cacheMetric.MType == COUNTER {
			newValue := *cacheMetric.Delta + *metric.Delta
			d.cache[metric.ID].Delta = &newValue
		}
		d.cache[metric.ID].Value = metric.Value

	}
	valuesQuery := make([]string, 0, len(metricSlice))
	valuesArgs := make([]interface{}, 0, len(metricSlice)*4)
	i := 1
	for _, metric := range d.cache {
		valuesQuery = append(valuesQuery, fmt.Sprintf("($%d, $%d, $%d, $%d)", i, i+1, i+2, i+3))
		i += 4
		valuesArgs = append(valuesArgs, metric.ID, metric.MType, metric.Value, metric.Delta)
	}
	insertBatchQuery := fmt.Sprintf(insertBatchQuery, strings.Join(valuesQuery, ", "))
	_, err := d.db.Exec(insertBatchQuery, valuesArgs...)
	if err != nil {
		return err
	}
	return nil
}
