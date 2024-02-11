package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHashSignMiddleWare(t *testing.T) {
	router := gin.Default()

	cfg := &config.Config{
		HashKey: "testkey",
	}

	router.Use(HashSignMiddleWare(cfg))

	router.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test response")
	})

	req, err := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte("test data")))
	if err != nil {
		t.Fatal(err)
	}
	hash := hmac.New(sha256.New, []byte(cfg.HashKey))
	_, err = hash.Write([]byte("test data"))
	if err != nil {
		t.Fatal(err)
	}
	sign := hex.EncodeToString(hash.Sum(nil))
	req.Header.Set("HashSHA256", sign)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	req, err = http.NewRequest("POST", "/test", bytes.NewBuffer([]byte("test data")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("HashSHA256", "invalidsignature")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, w.Code)
	}

	req, err = http.NewRequest("POST", "/test", bytes.NewBuffer([]byte("test data")))
	if err != nil {
		t.Fatal(err)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
}
