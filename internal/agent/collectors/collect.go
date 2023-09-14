package collectors

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/Kotletta-TT/MonoGo/internal/agent/utils"
	"log"
)

type collectorStorage interface {
	StoreMetrics(map[string]*entity.Value)
}

type Collector struct {
	repo               *collectorStorage
	registerCollectors []func(repo collectorStorage)
}

func NewCollector(repo collectorStorage) Collector {
	log.Println("Create collector")
	if repo == nil {
		panic("storage have nil pointer")
	}
	return Collector{
		repo:               &repo,
		registerCollectors: make([]func(repo collectorStorage), 0),
	}
}

func (c *Collector) RegisterCollectorMetricFunc(f func(repo collectorStorage)) {
	log.Printf("Register collectors func %s\n", utils.GetFunctionName(f))
	c.registerCollectors = append(c.registerCollectors, f)
}

func (c *Collector) Collect() {
	log.Printf("start collect, count registred func: %d\n", len(c.registerCollectors))
	for _, f := range c.registerCollectors {
		f(*c.repo)
	}
}
