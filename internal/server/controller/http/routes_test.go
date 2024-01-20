package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage/memory"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestListMetrics(t *testing.T) {
	tests := []struct {
		name string
		want []byte
		cfg  *config.Config
	}{
		{
			name: "Empty 200 OK",
			want: make([]byte, 0),
			cfg: &config.Config{
				RunServerAddr:   "localhost:8080",
				LogLevel:        "INFO",
				LogPath:         "",
				LogFile:         false,
				StoreInterval:   300,
				FileStoragePath: "",
				Restore:         false,
			},
		},
	}
	gin.SetMode(gin.TestMode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := memory.New(tt.cfg)
			r := gin.Default()
			r.GET("/", ListMetrics(repo))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			assert.True(t, bytes.Equal(tt.want, w.Body.Bytes()))
		})
	}
}

func TestGetMetric(t *testing.T) {
	repo := &RepositoryMock{}

	handler := GetMetric(repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)
	ctx, _ := gin.CreateTestContext(w)

	handler(ctx)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "", w.Body.String())

	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	handler(ctx)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	handler(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSetMetric(t *testing.T) {
	repo := &RepositoryMock{}

	handler := SetMetric(repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/test", handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)
	ctx, _ := gin.CreateTestContext(w)

	handler(ctx)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	handler(ctx)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(w)
	handler(ctx)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSetBatchJSONMetric(t *testing.T) {
	repo := &RepositoryMock{}

	handler := SetBatchJSONMetric(repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/test", handler)

	jsonData := `[{"metricType": "type1", "metric": "metric1", "value": 42}, {"metricType": "type2", "metric": "metric2", "value": 99}]`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	ctx, _ := gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", strings.NewReader("invalid json data"))
	router.ServeHTTP(w, req)
	ctx, _ = gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	router.ServeHTTP(w, req)
	ctx, _ = gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestSetJSONMetric(t *testing.T) {
	repo := &RepositoryMock{}

	handler := SetJSONMetric(repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/test", handler)

	jsonData := `{"metricType": "type1", "metric": "metric1", "value": 42}`
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	ctx, _ := gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "{\"error\":\"metric name is empty\"}", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", strings.NewReader("invalid json data"))
	router.ServeHTTP(w, req)
	ctx, _ = gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", strings.NewReader(jsonData))
	router.ServeHTTP(w, req)
	ctx, _ = gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetJSONMetric(t *testing.T) {
	repo := &RepositoryMock{}

	handler := GetJSONMetric(repo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/test", handler)

	jsonData := `{"metricType": "type1", "metric": "metric1", "value": 42}`
	req, _ := http.NewRequest("GET", "/test", strings.NewReader(jsonData))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	ctx, _ := gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"error\":\"metric name is empty\"}", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", strings.NewReader("invalid json data"))
	router.ServeHTTP(w, req)
	ctx, _ = gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", strings.NewReader(jsonData))
	router.ServeHTTP(w, req)
	ctx, _ = gin.CreateTestContext(w)
	assert.Panics(t, func() {
		handler(ctx)
	})
	assert.Equal(t, http.StatusNotFound, w.Code)
}
