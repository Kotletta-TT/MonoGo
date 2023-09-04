package usecase

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/entity"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"testing"
)

func TestCustomMetricsCollector(t *testing.T) {
	type args struct {
		repo       repository.Repository
		lenRepo    int
		iterate    int
		countValue int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Zero iterate collect",
			args: args{repo: repository.NewMemRepo(), lenRepo: 0, iterate: 0, countValue: 0},
		},
		{
			name: "Normal one iterate collect",
			args: args{repo: repository.NewMemRepo(), lenRepo: 2, iterate: 1, countValue: 1},
		},
		{
			name: "Normal two iterate collect",
			args: args{repo: repository.NewMemRepo(), lenRepo: 2, iterate: 2, countValue: 2},
		},
		{
			name: "Normal ten iterate collect",
			args: args{repo: repository.NewMemRepo(), lenRepo: 2, iterate: 10, countValue: 10},
		},
		{
			name: "Normal one hungered iterate collect",
			args: args{repo: repository.NewMemRepo(), lenRepo: 2, iterate: 100, countValue: 100},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.args.iterate; i++ {
				CustomMetricsCollector(tt.args.repo)
			}
			if len(tt.args.repo.GetAllMetrics()) != tt.args.lenRepo {
				t.Errorf("Expected len: %d got: %d", tt.args.lenRepo, len(tt.args.repo.GetAllMetrics()))
			}
			metric, err := tt.args.repo.GetMetric("PollCount")
			if err != nil && tt.args.iterate == 0 {
				return
			} else if err != nil && tt.args.iterate != 0 {
				t.Errorf("iterate not zero, but get metric err %s", err)
			}
			if metric.GetCounterValue() != int64(tt.args.countValue) {
				t.Errorf("Expected count value: %d got: %d", tt.args.countValue, metric.GetCounterValue())
			}
		})
	}
}

func TestRuntimeMetricsCollector(t *testing.T) {
	type args struct {
		repo        repository.Repository
		lenRepo     int
		excludeType entity.MetricKind
		iterate     int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Zero iterate collect",
			args: args{
				repo:        repository.NewMemRepo(),
				lenRepo:     0,
				excludeType: entity.KindCounter,
				iterate:     0,
			},
		},
		{
			name: "Normal one iterate collect",
			args: args{
				repo:        repository.NewMemRepo(),
				lenRepo:     27,
				excludeType: entity.KindCounter,
				iterate:     1,
			},
		},
		{
			name: "Normal two iterate collect",
			args: args{
				repo:        repository.NewMemRepo(),
				lenRepo:     27,
				excludeType: entity.KindCounter,
				iterate:     2,
			},
		},
		{
			name: "Normal ten iterate collect",
			args: args{
				repo:        repository.NewMemRepo(),
				lenRepo:     27,
				excludeType: entity.KindCounter,
				iterate:     10,
			},
		},
		{
			name: "Normal one hungered iterate collect",
			args: args{
				repo:        repository.NewMemRepo(),
				lenRepo:     27,
				excludeType: entity.KindCounter,
				iterate:     100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.args.iterate; i++ {
				RuntimeMetricsCollector(tt.args.repo)
			}
			all := tt.args.repo.GetAllMetrics()
			if len(all) != tt.args.lenRepo {
				t.Errorf("Expected len: %d got: %d", tt.args.lenRepo, len(tt.args.repo.GetAllMetrics()))
			}
			for _, m := range all {
				if m.GetMetricKind() == tt.args.excludeType {
					t.Errorf("Expected metric type: %d got: %d metric: %s", entity.KindGauge, tt.args.excludeType, m.Name)
				}
			}
		})
	}
}
