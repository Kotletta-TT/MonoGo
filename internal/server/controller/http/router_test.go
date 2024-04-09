package http

import (
	"context"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/common"
	log "github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestNewRouter(t *testing.T) {
	// Создаем тестовый репозиторий и конфиг

	cfg := &config.Config{
		HashKey: "testkey",
	}
	log.Init(cfg)
	repo := &RepositoryMock{} // Замените на вашу реализацию репозитория
	// Создаем роутер с использованием вашей функции
	router := NewRouter(repo, cfg)

	// Создаем тестовый запрос
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Записываем запрос в буфер
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем, что роутер включает в себя все необходимые middleware и обработчики
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}
