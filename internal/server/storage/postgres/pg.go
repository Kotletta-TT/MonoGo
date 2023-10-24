package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

const insertQueryMetric = `
INSERT INTO metrics (name, mtype, value, delta)
VALUES ($1, $2, $3, $4)
ON CONFLICT (name)
DO UPDATE SET
	value = EXCLUDED.value,
	delta = metrics.delta + EXCLUDED.delta
RETURNING name, mtype, value, delta
`
const insertQueryMetrics = `
	INSERT INTO metrics (name, mtype, value, delta)
	VALUES %s
	ON CONFLICT (name)
	DO UPDATE SET 
		value = EXCLUDED.value,
		delta = metrics.delta + EXCLUDED.delta
	RETURNING name, mtype, value, delta
`

const selectQueryMetric = `SELECT name, mtype, value, delta FROM metrics WHERE name = $1 AND mtype = $2`
const selectQueryMetrics = `SELECT * FROM metrics`

const insertQueryMetricStart = `INSERT INTO metrics (name, mtype, value, delta) VALUES `
const insertQueryMetricEnd = ` ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta RETURNING name, mtype, value, delta`

const createTableQuery = `
	CREATE TABLE 
	IF NOT EXISTS metrics (
	name VARCHAR(255) NOT NULL,
	mtype VARCHAR(10) NOT NULL,
	value DOUBLE PRECISION,
	delta BIGINT,
	PRIMARY KEY (name));`

type Database struct {
	cfg    *config.Config
	pgPool *pgxpool.Pool
	ctx    context.Context
}

func New(cfg *config.Config) (*Database, error) {
	ctx := context.Background()
	pgConn, err := pgxpool.New(ctx, cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}
	err = pgConn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	pgConn.Exec(ctx, createTableQuery)
	return &Database{
		pgPool: pgConn,
		cfg:    cfg,
		ctx:    ctx,
	}, nil
}

func (d *Database) WrapPgError(pgFunc func() error) error {
	var err error
	attempt := 3
	timeoutRerun := 1
	for i := 0; i <= attempt; i++ {
		err = pgFunc()
		if err == nil {
			break
		}
		var pgError *pgx.PgError
		if errors.As(err, &pgError) && pgerrcode.IsConnectionException(pgError.Code) {
			logger.Errorf("Connect to database error: %s attempt %d", err, i+1)
			time.Sleep(time.Duration(time.Second * time.Duration(timeoutRerun)))
			timeoutRerun += 2
			continue
		}
		return err
	}
	return err
}

func (d *Database) StoreMetric(metric *common.Metrics) error {
	return d.WrapPgError(func() error {
		row := d.pgPool.QueryRow(d.ctx, insertQueryMetric, metric.ID, metric.MType, metric.Value, metric.Delta)
		err := row.Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *Database) makeInsertBatchQueryValues(metrics []*common.Metrics) (string, []interface{}) {
	valuesQuery := make([]string, 0, len(metrics))
	valuesArgs := make([]interface{}, 0, len(metrics)*4)
	i := 1
	for _, metric := range metrics {
		valuesQuery = append(valuesQuery, fmt.Sprintf("($%d, $%d, $%d, $%d)", i, i+1, i+2, i+3))
		i += 4
		valuesArgs = append(valuesArgs, metric.ID, metric.MType, metric.Value, metric.Delta)
	}
	return insertQueryMetricStart + strings.Join(valuesQuery, ", ") + insertQueryMetricEnd, valuesArgs
}

func (d *Database) StoreBatchMetric(metrics []*common.Metrics) ([]*common.Metrics, error) {
	var resultMetrics []*common.Metrics
	err := d.WrapPgError(func() error {
		var errs error
		if len(metrics) == 0 {
			errs = fmt.Errorf("metrics can't be empty")
			return errs
		}
		compiledQuery, compiledValues := d.makeInsertBatchQueryValues(metrics)
		rows, errs := d.pgPool.Query(d.ctx, compiledQuery, compiledValues...)
		if errs != nil {
			return errs
		}
		defer rows.Close()
		resultMetrics = make([]*common.Metrics, 0, len(metrics))
		for rows.Next() {
			m := new(common.Metrics)
			errs = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
			if errs != nil {
				return errs
			}
			if rows.Err() != nil {
				return rows.Err()
			}
			resultMetrics = append(resultMetrics, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return resultMetrics, nil
}

func (d *Database) GetMetric(metric *common.Metrics) error {
	return d.WrapPgError(func() error {
		err := d.pgPool.QueryRow(d.ctx, selectQueryMetric, metric.ID, metric.MType).Scan(&metric.ID, &metric.MType, &metric.Value, &metric.Delta)
		if err != nil {
			return err
		}
		return nil
	})
}

func (d *Database) GetListMetrics() ([]*common.Metrics, error) {
	var metrics []*common.Metrics
	err := d.WrapPgError(func() error {
		var errs error
		rows, errs := d.pgPool.Query(d.ctx, selectQueryMetrics)
		if errs != nil {
			return errs
		}
		defer rows.Close()
		metrics = make([]*common.Metrics, 0, 10)
		for rows.Next() {
			m := new(common.Metrics)
			errs = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
			if errs != nil {
				return errs
			}
			if rows.Err() != nil {
				return rows.Err()
			}
			metrics = append(metrics, m)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (d *Database) HealthCheck(ctx context.Context) error {
	return d.pgPool.Ping(ctx)
}

func (d *Database) Close() {
	d.pgPool.Close()
}
