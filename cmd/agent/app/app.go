// Package app implements the main process of the application.
package app

import (
	"log"
	"time"

	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
	"github.com/Kotletta-TT/MonoGo/internal/agent/collectors"
	"github.com/Kotletta-TT/MonoGo/internal/agent/sender"
	"github.com/Kotletta-TT/MonoGo/internal/agent/storage"
)

// Run runs the main process of the application.
//
// It takes a pointer to a Config struct as a parameter.
// It does not return anything.
func Run(cfg *config.Config) {
	repo := storage.New()
	collector := collectors.NewCollector(repo)
	collector.RegisterCollectorMetricFunc(collectors.RuntimeMetricsCollector)
	collector.RegisterCollectorMetricFunc(collectors.CustomMetricsCollector)
	collector.RegisterCollectorMetricFunc(collectors.SystemStatsCollector)
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
