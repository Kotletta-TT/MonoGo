package http

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCompressMiddleware(t *testing.T) {
	router := gin.Default()

	router.Use(CompressMiddleware())

	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Test response")
	})

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	respBody, err := gzip.NewReader(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer respBody.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(respBody)
	if err != nil {
		t.Fatal(err)
	}

	expectedResponse := "Test response"
	if buf.String() != expectedResponse {
		t.Errorf("Expected response %q, but got %q", expectedResponse, buf.String())
	}

	req, err = http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	if w.Body.String() != expectedResponse {
		t.Errorf("Expected response %q, but got %q", expectedResponse, w.Body.String())
	}
}
