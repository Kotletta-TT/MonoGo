package http

import (
	"context"
	"github.com/Kotletta-TT/MonoGo/internal/common"
)

type RepositoryMock struct{}

func (r RepositoryMock) StoreMetric(metric *common.Metrics) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) StoreBatchMetric(metrics []*common.Metrics) ([]*common.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) GetMetric(metric *common.Metrics) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) GetListMetrics() ([]*common.Metrics, error) {
	return []*common.Metrics{}, nil
}

func (r RepositoryMock) HealthCheck(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r RepositoryMock) Close() {
	//TODO implement me
	panic("implement me")
}
