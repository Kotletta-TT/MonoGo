package app

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/infrastructure/repository"
	"github.com/Kotletta-TT/MonoGo/internal/agent/usecase"
	"log"
	"time"
)

var poolInterval = time.Second * 2
var reportInterval = time.Second * 10

func Run() {
	repo := repository.New()
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
			go sender.Send()
		}
	}
}
