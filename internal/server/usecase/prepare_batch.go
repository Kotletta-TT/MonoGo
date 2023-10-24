package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/common"
)

func PrepareBatchStore(metrics []*common.Metrics) []*common.Metrics {
	sortMapMetrics := make(map[string]*common.Metrics)
	for _, m := range metrics {
		newMetric, ok := sortMapMetrics[m.ID]
		if !ok {
			sortMapMetrics[m.ID] = m
			continue
		}
		switch m.MType {
		case GAUGE:
			newMetric.Value = m.Value
		case COUNTER:
			newMetric.Delta = common.SumDelta(newMetric.Delta, m.Delta)
		}
	}
	newArr := make([]*common.Metrics, 0, len(sortMapMetrics))
	for _, sortedM := range sortMapMetrics {
		newArr = append(newArr, sortedM)
	}
	return newArr
}
