package app

import (
	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
	"github.com/Kotletta-TT/MonoGo/internal/agent/collectors"
	"github.com/Kotletta-TT/MonoGo/internal/agent/sender"
	"github.com/Kotletta-TT/MonoGo/internal/agent/storage"
	"log"
	"time"
)

func Run(cfg *config.Config) {
	repo := storage.New()
	collector := collectors.NewCollector(repo)
	collector.RegisterCollectorMetricFunc(collectors.RuntimeMetricsCollector)
	collector.RegisterCollectorMetricFunc(collectors.CustomMetricsCollector)
	httpSender := sender.NewHTTPSender(repo, cfg)
	pollTic := time.NewTicker(time.Second * time.Duration(cfg.PollInterval))
	reportTic := time.NewTicker(time.Second * time.Duration(cfg.ReportInterval))
	log.Println("Start work")
	for {
		select {
		case <-pollTic.C:
			go collector.Collect()
		case <-reportTic.C:
			go httpSender.Send()
		}
	}
}
