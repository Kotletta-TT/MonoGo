package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/utils"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"log"
)

type Collector struct {
	repo               *repository.Repository
	registerCollectors []func(repo repository.Repository)
}

func NewCollector(repo *repository.Repository) Collector {
	log.Println("Create collector")
	if repo == nil {
		panic("repository have nil pointer")
	}
	return Collector{
		repo:               repo,
		registerCollectors: make([]func(repo repository.Repository), 0),
	}
}

func (c *Collector) RegisterCollectorMetricFunc(f func(repo repository.Repository)) {
	log.Printf("Register collectors func %s\n", utils.GetFunctionName(f))
	c.registerCollectors = append(c.registerCollectors, f)
}

func (c *Collector) Collect() {
	log.Printf("start collect, count registred func: %d\n", len(c.registerCollectors))
	for _, f := range c.registerCollectors {
		f(*c.repo)
	}
}
