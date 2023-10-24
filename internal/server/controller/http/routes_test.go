package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
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
