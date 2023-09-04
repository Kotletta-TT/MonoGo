package app

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/usecase"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"log"
	"time"
)

var poolInterval = time.Second * 2
var reportInterval = time.Second * 10

func Run() {
	repo := repository.NewMemRepo()
	collector := usecase.NewCollector(&repo)
	collector.RegisterCollectorMetricFunc(usecase.RuntimeMetricsCollector)
	collector.RegisterCollectorMetricFunc(usecase.CustomMetricsCollector)
	sender := usecase.NewHTTPSender(repo)
	poolTic := time.NewTicker(poolInterval)
	reportTic := time.NewTicker(reportInterval)
	log.Println("Start work")
	for {
		select {
		case <-poolTic.C:
			go collector.Collect()
		case <-reportTic.C:
			go repo.PrintAllMetrics()
			go sender.Send()
		}
	}
}
