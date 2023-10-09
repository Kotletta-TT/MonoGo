package memory

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"maps"
	"os"
	"sync"
	"time"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/shared"
	"github.com/mailru/easyjson"
)

const (
	GAUGE   = "gauge"
	COUNTER = "counter"
)

type MemRepository struct {
	mu      sync.Mutex
	storage map[string]*shared.Metrics
	cfg     *config.Config
}

func New(cfg *config.Config) *MemRepository {
	store := make(map[string]*shared.Metrics)
	m := &MemRepository{mu: sync.Mutex{}, cfg: cfg}
	if cfg.Restore {
		logger.Infof("Attempt to restore from file: %s", cfg.FileStoragePath)
		loadBackup, err := m.LoadFromFile()
		if err == nil {
			store = loadBackup
		} else {
			logger.Error(err)
		}
	}
	m.storage = store
	return m
}

func (m *MemRepository) StoreGaugeMetric(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage[name] = &shared.Metrics{
		ID:    name,
		MType: GAUGE,
		Value: &value,
	}
	if m.cfg.StoreInterval == 0 {
		m.store()
	}
}

func (m *MemRepository) StoreCounterMetric(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.storage[name]
	if !ok {
		m.storage[name] = &shared.Metrics{
			ID:    name,
			MType: COUNTER,
			Delta: &value,
		}
		return
	}
	valInt := *val.Delta + value
	m.storage[name].Delta = &valInt
	if m.cfg.StoreInterval == 0 {
		m.store()
	}
}

func (m *MemRepository) GetGaugeMetric(name string) (float64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metric, ok := m.storage[name]
	if !ok {
		return 0.0, fmt.Errorf("metric not found")
	}
	return *metric.Value, nil
}

func (m *MemRepository) GetCounterMetric(name string) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	metric, ok := m.storage[name]
	if !ok {
		return 0.0, fmt.Errorf("metric not found")
	}
	return *metric.Delta, nil
}

func (m *MemRepository) GetAllMetrics() map[string]*shared.Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.storage
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

func (m *MemRepository) LoadFromFile() (map[string]*shared.Metrics, error) {
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
	metrics := make(map[string]*shared.Metrics)
	for scn.Scan() {
		line := scn.Bytes()
		if len(line) == 0 {
			continue
		}
		metric := shared.NewMetrics()
		if err := easyjson.Unmarshal(line, metric); err != nil {
			return nil, err
		}
		metrics[metric.ID] = metric
	}
	return metrics, err
}

func (m *MemRepository) Close() {}

func (m *MemRepository) HealthCheck(ctx context.Context) error {
	return errors.New("No connect to database")
}
