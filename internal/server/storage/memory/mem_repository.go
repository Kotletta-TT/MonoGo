package memory

import (
	"bufio"
	"context"
	"errors"
	"maps"
	"os"
	"sync"
	"time"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/mailru/easyjson"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

type MemRepository struct {
	mu      sync.Mutex
	storage map[string]*common.Metrics
	cfg     *config.Config
}

// New creates a new MemRepository instance.
//
// It takes a pointer to a config.Config struct as its parameter.
// It returns a pointer to a MemRepository struct.
func New(cfg *config.Config) *MemRepository {
	store := make(map[string]*common.Metrics)
	m := &MemRepository{mu: sync.Mutex{}, cfg: cfg}
	if cfg.Restore {
		logger.Infof("Attempt to restore from file: %s", cfg.FileStoragePath)
		loadBackup, err := m.loadFromFile()
		if err == nil {
			store = loadBackup
		} else {
			logger.Error(err)
		}
	}
	m.storage = store
	return m
}

func (m *MemRepository) createOrUpdateMetric(metric *common.Metrics) {
	val, ok := m.storage[metric.ID]
	if !ok {
		m.storage[metric.ID] = metric
		return
	}
	switch metric.MType {
	case GAUGE:
		*m.storage[metric.ID] = *metric
	case COUNTER:
		m.storage[metric.ID].Delta = common.SumDelta(val.Delta, metric.Delta)
		*metric.Delta = *m.storage[metric.ID].Delta
	}
}

func (m *MemRepository) StoreMetric(metric *common.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createOrUpdateMetric(metric)
	if m.cfg.StoreInterval == 0 {
		m.store()
	}
	return nil
}

func (m *MemRepository) StoreBatchMetric(metrics []*common.Metrics) ([]*common.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, metric := range metrics {
		m.createOrUpdateMetric(metric)
	}
	if m.cfg.StoreInterval == 0 {
		m.store()
	}
	return metrics, nil
}

func (m *MemRepository) GetMetric(metric *common.Metrics) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.storage[metric.ID]
	if !ok {
		return errors.New("metric not found")
	}
	*metric = *val
	return nil
}

func (m *MemRepository) GetListMetrics() ([]*common.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sliceStorage := make([]*common.Metrics, 0, len(m.storage))
	for _, v := range m.storage {
		sliceStorage = append(sliceStorage, v)
	}
	return sliceStorage, nil
}

func (m *MemRepository) store() {
	file, err := os.OpenFile(m.cfg.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		if err = file.Close(); err != nil {
			logger.Error(err)
		}
	}()
	m.mu.Lock()
	cloneStore := maps.Clone(m.storage)
	m.mu.Unlock()
	for _, v := range cloneStore {
		_, err = file.WriteString(v.String() + "\n")
		if err != nil {
			logger.Error(err)

			continue
		}
	}
}

func (m *MemRepository) Stash() {
	if m.cfg.StoreInterval == 0 {
		return
	}
	stashTic := time.NewTicker(time.Second * time.Duration(m.cfg.StoreInterval))
	for range stashTic.C {
		m.store()
	}
}

func (m *MemRepository) loadFromFile() (map[string]*common.Metrics, error) {
	file, err := os.Open(m.cfg.FileStoragePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Infof("Error closing file: %s", err)
		}
	}()
	scn := bufio.NewScanner(file)
	metrics := make(map[string]*common.Metrics)
	for scn.Scan() {
		line := scn.Bytes()
		if len(line) == 0 {
			continue
		}
		metric := new(common.Metrics)
		if err := easyjson.Unmarshal(line, metric); err != nil {
			return nil, err
		}
		metrics[metric.ID] = metric
	}
	return metrics, err
}

func (m *MemRepository) Close() {}

func (m *MemRepository) HealthCheck(ctx context.Context) error {
	return errors.New("no connect to database")
}
