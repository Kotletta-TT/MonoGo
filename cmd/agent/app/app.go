package app

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/infrastructure/repository"
	"github.com/Kotletta-TT/MonoGo/internal/agent/usecase"
	"log"
	"time"
)

func Run() {
	parseFlags()
	repo := repository.New()
	collector := usecase.NewCollector(&repo)
	collector.RegisterCollectorMetricFunc(usecase.RuntimeMetricsCollector)
	collector.RegisterCollectorMetricFunc(usecase.CustomMetricsCollector)
	sender := usecase.NewHTTPSender(repo, flagServerAddr)
	poolTic := time.NewTicker(time.Second * time.Duration(flagPoolInterval))
	reportTic := time.NewTicker(time.Second * time.Duration(flagReportInterval))
	log.Println("Start work")
	for {
		select {
		case <-poolTic.C:
			go collector.Collect()
		case <-reportTic.C:
			go sender.Send()
		}
	}
}
