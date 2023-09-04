package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/entity"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"log"
	"math/rand"
	"reflect"
	"runtime"
)

var runtimeMetrics = map[string]entity.MetricKind{
	"Alloc":         entity.KindGauge,
	"BuckHashSys":   entity.KindGauge,
	"Frees":         entity.KindGauge,
	"GCCPUFraction": entity.KindGauge,
	"GCSys":         entity.KindGauge,
	"HeapAlloc":     entity.KindGauge,
	"HeapIdle":      entity.KindGauge,
	"HeapInuse":     entity.KindGauge,
	"HeapObjects":   entity.KindGauge,
	"HeapReleased":  entity.KindGauge,
	"HeapSys":       entity.KindGauge,
	"LastGC":        entity.KindGauge,
	"Lookups":       entity.KindGauge,
	"MCacheInuse":   entity.KindGauge,
	"MCacheSys":     entity.KindGauge,
	"MSpanInuse":    entity.KindGauge,
	"MSpanSys":      entity.KindGauge,
	"Mallocs":       entity.KindGauge,
	"NextGC":        entity.KindGauge,
	"NumForcedGC":   entity.KindGauge,
	"NumGC":         entity.KindGauge,
	"OtherSys":      entity.KindGauge,
	"PauseTotalNs":  entity.KindGauge,
	"StackInuse":    entity.KindGauge,
	"StackSys":      entity.KindGauge,
	"Sys":           entity.KindGauge,
	"TotalAlloc":    entity.KindGauge,
}

var customMetrics = map[string]entity.MetricKind{
	"PollCount":   entity.KindCounter,
	"RandomValue": entity.KindGauge,
}

func RuntimeMetricsCollector(repo repository.Repository) {
	log.Println("start runtime metrics collector")
	rawMetrics := runtime.MemStats{}
	runtime.ReadMemStats(&rawMetrics)
	val := reflect.ValueOf(rawMetrics)
	typeOfMemStats := val.Type()
	for i := 0; i < val.NumField(); i++ {
		name := typeOfMemStats.Field(i).Name
		_, ok := runtimeMetrics[name]
		if !ok {
			continue
		}
		switch runtimeMetrics[name] {
		case entity.KindCounter:
			repo.StoreCounterMetric(name, val.Field(i).Interface())
		case entity.KindGauge:
			repo.StoreGaugeMetric(name, val.Field(i).Interface())
		default:
			//TODO обработка?
			continue
		}
	}
}

func CustomMetricsCollector(repo repository.Repository) {
	pollCount := 1
	randValue := rand.Float64()
	log.Println("start custom metrics collector")
	for k, v := range customMetrics {
		switch v {
		case entity.KindCounter:
			repo.StoreCounterMetric(k, pollCount)
		case entity.KindGauge:
			repo.StoreGaugeMetric(k, randValue)
		default:
			continue
		}
	}
}
