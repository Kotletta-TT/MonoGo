package usecase

import (
	"encoding/json"
	"github.com/Kotletta-TT/MonoGo/internal/agent/entity"
	"github.com/Kotletta-TT/MonoGo/internal/agent/infrastructure/repository"
	"log"
	"math/rand"
	"runtime"
)

func RuntimeMetricsCollector(repo repository.AgentRepository) {
	log.Println("start runtime metrics collector")
	var finalFilterMap map[string]interface{}
	runtimeMetrics := runtime.MemStats{}
	runtime.ReadMemStats(&runtimeMetrics)
	bRuntimeMetrics, err := json.Marshal(runtimeMetrics)
	if err != nil {
		panic(err)
	}
	neededMetrics := entity.RuntimeMetrics{}
	if err := json.Unmarshal(bRuntimeMetrics, &neededMetrics); err != nil {
		panic(err)
	}
	bNeededMetrics, err := json.Marshal(neededMetrics)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bNeededMetrics, &finalFilterMap); err != nil {
		panic(err)
	}
	repo.StoreMetrics(finalFilterMap)
}

func CustomMetricsCollector(repo repository.AgentRepository) {
	log.Println("start custom metrics collector")
	pollCount := int64(1)
	randValue := rand.Float64()
	repo.StoreMetrics(map[string]interface{}{"PoolCount": pollCount, "RandomValue": randValue})
}
