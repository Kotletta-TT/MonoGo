package http

import (
	"bytes"
	"fmt"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	log "github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestResponseLogging(t *testing.T) {
	cfg := config.NewConfig()
	log.Init(cfg)
	router := gin.Default()

	router.Use(RequestResponseLogging())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test response")
	})

	buf := new(bytes.Buffer)
	gin.DefaultWriter = buf

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	loggedOutput := buf.String()
	expectedLog := fmt.Sprintf("%s %s %d %s", "127.0.0.1", "GET", http.StatusOK, "/test")
	if !contains(loggedOutput, expectedLog) {
		t.Logf("Expected log to contain %q, but got %q", expectedLog, loggedOutput)
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
