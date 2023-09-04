package repository

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/internal/agent/utils"
	"github.com/Kotletta-TT/MonoGo/internal/server/entity"
	"log"
	"sync"
)

type Repository interface {
	GetMetric(name string) (*entity.CustomMetric, error)
	StoreGaugeMetric(key string, value interface{})
	StoreCounterMetric(key string, value interface{})
	GetAllMetrics() map[string]*entity.CustomMetric
	PrintAllMetrics()
}

type MemRepo struct {
	mu   *sync.Mutex
	repo map[string]*entity.CustomMetric
}

func NewMemRepo() Repository {
	log.Println("Create Mem Repository")
	return &MemRepo{
		mu:   &sync.Mutex{},
		repo: make(map[string]*entity.CustomMetric),
	}
}

func (mr *MemRepo) GetMetric(name string) (*entity.CustomMetric, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	m, ok := mr.repo[name]
	if !ok {
		return nil, fmt.Errorf("metric %s not found", name)
	}
	return m, nil
}

func (mr *MemRepo) StoreGaugeMetric(key string, value interface{}) {
	floatValue := utils.GetFloat64(value)
	mr.mu.Lock()
	defer mr.mu.Unlock()
	if _, ok := mr.repo[key]; !ok {
		mr.repo[key] = entity.NewCustomGaugeMetric(key, floatValue)
		return
	}
	mr.repo[key].UpdateGageValue(floatValue)
}

func (mr *MemRepo) StoreCounterMetric(key string, value interface{}) {
	intValue := utils.GetInt64(value)
	mr.mu.Lock()
	defer mr.mu.Unlock()
	if _, ok := mr.repo[key]; !ok {
		mr.repo[key] = entity.NewCustomCounterMetric(key, intValue)
		return
	}
	mr.repo[key].UpdateCounterValue(intValue)
}

func (mr *MemRepo) PrintAllMetrics() {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	log.Printf("Stored metrics in repo: %d\n", len(mr.repo))
	for _, metric := range mr.repo {
		switch metric.GetMetricKind() {
		case entity.KindGauge:
			log.Printf("Name:%s Value:%v Kind: Gauge\n", metric.Name, metric.GetGaugeValue())
		case entity.KindCounter:
			log.Printf("Name:%s Value:%v Kind: Gauge\n", metric.Name, metric.GetCounterValue())
		}
	}
}

func (mr *MemRepo) GetAllMetrics() map[string]*entity.CustomMetric {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	copyRepo := make(map[string]*entity.CustomMetric, len(mr.repo))
	for k, v := range mr.repo {
		copyRepo[k] = v
	}
	return copyRepo
}
