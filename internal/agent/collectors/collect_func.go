package collectors

import (
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"log"
	"math"
	"math/rand"
	"runtime"
)

func RuntimeMetricsCollector(repo collectorStorage) {
	log.Println("start runtime metrics collector")
	metricMap := make(map[string]*entity.Value)
	runtimeMetrics := runtime.MemStats{}
	runtime.ReadMemStats(&runtimeMetrics)
	metricMap["Alloc"] = entity.NewValueFromUint64(runtimeMetrics.Alloc, entity.KindGauge)
	metricMap["BuckHashSys"] = entity.NewValueFromUint64(runtimeMetrics.BuckHashSys, entity.KindGauge)
	metricMap["Frees"] = entity.NewValueFromUint64(runtimeMetrics.Frees, entity.KindGauge)
	metricMap["GCCPUFraction"] = entity.NewValueFromFloat64(runtimeMetrics.GCCPUFraction, entity.KindGauge)
	metricMap["GCSys"] = entity.NewValueFromUint64(runtimeMetrics.GCSys, entity.KindGauge)
	metricMap["HeapAlloc"] = entity.NewValueFromUint64(runtimeMetrics.HeapAlloc, entity.KindGauge)
	metricMap["HeapIdle"] = entity.NewValueFromUint64(runtimeMetrics.HeapIdle, entity.KindGauge)
	metricMap["HeapInuse"] = entity.NewValueFromUint64(runtimeMetrics.HeapInuse, entity.KindGauge)
	metricMap["HeapObjects"] = entity.NewValueFromUint64(runtimeMetrics.HeapObjects, entity.KindGauge)
	metricMap["HeapReleased"] = entity.NewValueFromUint64(runtimeMetrics.HeapReleased, entity.KindGauge)
	metricMap["HeapSys"] = entity.NewValueFromUint64(runtimeMetrics.HeapSys, entity.KindGauge)
	metricMap["LastGC"] = entity.NewValueFromUint64(runtimeMetrics.LastGC, entity.KindGauge)
	metricMap["Lookups"] = entity.NewValueFromUint64(runtimeMetrics.Lookups, entity.KindGauge)
	metricMap["MCacheInuse"] = entity.NewValueFromUint64(runtimeMetrics.MCacheInuse, entity.KindGauge)
	metricMap["MCacheSys"] = entity.NewValueFromUint64(runtimeMetrics.MCacheSys, entity.KindGauge)
	metricMap["MSpanInuse"] = entity.NewValueFromUint64(runtimeMetrics.MSpanInuse, entity.KindGauge)
	metricMap["MSpanSys"] = entity.NewValueFromUint64(runtimeMetrics.MSpanSys, entity.KindGauge)
	metricMap["Mallocs"] = entity.NewValueFromUint64(runtimeMetrics.Mallocs, entity.KindGauge)
	metricMap["NextGC"] = entity.NewValueFromUint64(runtimeMetrics.NextGC, entity.KindGauge)
	metricMap["NumForcedGC"] = entity.NewValueFromUint32(runtimeMetrics.NumForcedGC, entity.KindGauge)
	metricMap["NumGC"] = entity.NewValueFromUint32(runtimeMetrics.NumGC, entity.KindGauge)
	metricMap["OtherSys"] = entity.NewValueFromUint64(runtimeMetrics.OtherSys, entity.KindGauge)
	metricMap["PauseTotalNs"] = entity.NewValueFromUint64(runtimeMetrics.PauseTotalNs, entity.KindGauge)
	metricMap["StackInuse"] = entity.NewValueFromUint64(runtimeMetrics.StackInuse, entity.KindGauge)
	metricMap["StackSys"] = entity.NewValueFromUint64(runtimeMetrics.StackSys, entity.KindGauge)
	metricMap["Sys"] = entity.NewValueFromUint64(runtimeMetrics.Sys, entity.KindGauge)
	metricMap["TotalAlloc"] = entity.NewValueFromUint64(runtimeMetrics.TotalAlloc, entity.KindGauge)
	repo.StoreMetrics(metricMap)
}

func CustomMetricsCollector(repo collectorStorage) {
	log.Println("start custom metrics collector")
	pollCount := &entity.Value{Metric: uint64(int64(1)), Kind: entity.KindCounter}
	randValue := &entity.Value{Metric: math.Float64bits(rand.Float64()), Kind: entity.KindGauge}
	repo.StoreMetrics(map[string]*entity.Value{"PollCount": pollCount, "RandomValue": randValue})
}
