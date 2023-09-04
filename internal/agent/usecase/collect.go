package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/infrastructure/repository"
	"github.com/Kotletta-TT/MonoGo/internal/agent/utils"
	"log"
)

type Collector struct {
	repo               *repository.AgentRepository
	registerCollectors []func(repo repository.AgentRepository)
}

func NewCollector(repo *repository.AgentRepository) Collector {
	log.Println("Create collector")
	if repo == nil {
		panic("repository have nil pointer")
	}
	return Collector{
		repo:               repo,
		registerCollectors: make([]func(repo repository.AgentRepository), 0),
	}
}

func (c *Collector) RegisterCollectorMetricFunc(f func(repo repository.AgentRepository)) {
	log.Printf("Register collectors func %s\n", utils.GetFunctionName(f))
	c.registerCollectors = append(c.registerCollectors, f)
}

func (c *Collector) Collect() {
	log.Printf("start collect, count registred func: %d\n", len(c.registerCollectors))
	for _, f := range c.registerCollectors {
		f(*c.repo)
	}
}
