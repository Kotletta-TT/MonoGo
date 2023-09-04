package http

import (
	"bytes"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListMetrics(t *testing.T) {
	tests := []struct {
		name string
		want []byte
	}{
		{
			name: "Empty 200 OK",
			want: make([]byte, 0),
		},
	}
	gin.SetMode(gin.TestMode)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repository.New()
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
