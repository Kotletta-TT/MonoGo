package usecase

import (
	"testing"

	"github.com/Kotletta-TT/MonoGo/internal/common"
	"github.com/stretchr/testify/assert"
)

func TestPrepareBatchStore(t *testing.T) {
	testCases := []struct {
		desc string
		src  []*common.Metrics
		exp  []*common.Metrics
	}{
		{
			desc: "One gauge element",
			src: []*common.Metrics{
				FabricGaugeMetric("metric1", 1.1),
			},
			exp: []*common.Metrics{
				FabricGaugeMetric("metric1", 1.1),
			},
		},
		{
			desc: "Two gauge duplicate elements",
			src: []*common.Metrics{
				FabricGaugeMetric("metric1", 1.1),
				FabricGaugeMetric("metric1", 2.2),
			},
			exp: []*common.Metrics{
				FabricGaugeMetric("metric1", 2.2),
			},
		},
		{
			desc: "One counter element",
			src: []*common.Metrics{
				FabricCounterMetric("metric1", 1),
			},
			exp: []*common.Metrics{
				FabricCounterMetric("metric1", 1),
			},
		},
		{
			desc: "Two counter duplicate element",
			src: []*common.Metrics{
				FabricCounterMetric("metric1", 1),
				FabricCounterMetric("metric1", 2),
			},
			exp: []*common.Metrics{
				FabricCounterMetric("metric1", 3),
			},
		},
		{
			desc: "Zero elements",
			src:  []*common.Metrics{},
			exp:  []*common.Metrics{},
		},
		{
			desc: "Counter1 Gauge1 Counter3 Counter1 Gauge1 Gauge1",
			src: []*common.Metrics{
				FabricCounterMetric("Counter1", 1),
				FabricGaugeMetric("Gauge1", 1.1),
				FabricCounterMetric("Counter3", 3),
				FabricCounterMetric("Counter1", 1),
				FabricGaugeMetric("Gauge1", 2.2),
				FabricGaugeMetric("Gauge1", 1.1),
			},
			exp: []*common.Metrics{
				FabricCounterMetric("Counter1", 2),
				FabricGaugeMetric("Gauge1", 1.1),
				FabricCounterMetric("Counter3", 3),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := PrepareBatchStore(tC.src)
			assert.Equal(t, len(tC.exp), len(res))
			for _, resValue := range res {
				for _, expValue := range tC.exp {
					if resValue.ID == expValue.ID {
						assert.Equal(t, expValue, resValue)
					}
				}
			}
		})
	}
}
