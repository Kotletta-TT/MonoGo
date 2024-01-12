// Package collectors implements the Collector object.
package collectors

import (
	"log"

	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/Kotletta-TT/MonoGo/internal/agent/utils"
)

type collectorStorage interface {
	StoreMetrics(map[string]*entity.Value)
}

// Collector collects metrics and stores them in the given repository.
//
// repo: the collector storage repository where the metrics will be stored.
// registerCollectors: a slice of functions that collect metrics and store them in the repository.
type Collector struct {
	repo               *collectorStorage
	registerCollectors []func(repo collectorStorage)
}

// NewCollector creates a new collector.
//
// It takes a repo collectorStorage parameter and returns a Collector.
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

// RegisterCollectorMetricFunc registers a collector metric function.
//
// The function takes a repo of collectorStorage as a parameter.
// There is no return type for this function.
func (c *Collector) RegisterCollectorMetricFunc(f func(repo collectorStorage)) {
	log.Printf("Register collectors func %s\n", utils.GetFunctionName(f))
	c.registerCollectors = append(c.registerCollectors, f)
}

// Collect collects data from registered collectors.
//
// It iterates over the registered collectors and calls each one,
// passing the repository as an argument.
func (c *Collector) Collect() {
	log.Printf("start collect, count registred func: %d\n", len(c.registerCollectors))
	for _, f := range c.registerCollectors {
		f(*c.repo)
	}
}
