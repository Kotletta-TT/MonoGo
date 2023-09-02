package repository

import (
	"github.com/Kotletta-TT/MonoGo/internal/entity"
)

type Repository interface {
	GetMetricByName(name string) (entity.Metric, error)
	StoreMetric(metric entity.Metric)
}

type MemRepo struct {
	repo map[string]entity.Metric
}

func NewMemRepo() MemRepo {
	return MemRepo{
		repo: make(map[string]entity.Metric),
	}
}

func (mr MemRepo) GetMetricByName(name string) (entity.Metric, error) {
	//TODO обработка отсутсвия метрики
	return mr.repo[name], nil
}

func (mr MemRepo) StoreMetric(metric entity.Metric) {
	mr.repo[metric.Name] = metric
}
