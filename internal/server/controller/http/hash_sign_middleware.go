package http

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/gin-gonic/gin"
)

type HMACResponseWriter struct {
	gin.ResponseWriter
	hashKey string
}

func (w HMACResponseWriter) Write(data []byte) (int, error) {
	h := hmac.New(sha256.New, []byte(w.hashKey))
	n, err := h.Write(data)
	if err != nil {
		return n, err
	}
	sign := hex.EncodeToString(h.Sum(nil))
	w.Header().Set("HashSHA256", sign)
	return w.ResponseWriter.Write(data)
}

func HashSignMiddleWare(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if cfg.HashKey != "" && ctx.GetHeader("HashSHA256") != "" {
			buf := bytes.NewBuffer(make([]byte, 0, ctx.Request.ContentLength))
			_, err := io.Copy(buf, ctx.Request.Body)
			if err != nil {
				logger.Error(err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			h := hmac.New(sha256.New, []byte(cfg.HashKey))
			_, err = h.Write(buf.Bytes())
			if err != nil {
				logger.Error(err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			calcSign := hex.EncodeToString(h.Sum(nil))
			if calcSign != ctx.GetHeader("HashSHA256") {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid sign"})
				return
			}
			ctx.Request.Body = io.NopCloser(buf)
			hmacWriter := &HMACResponseWriter{ResponseWriter: ctx.Writer, hashKey: cfg.HashKey}
			ctx.Writer = hmacWriter
		}
		ctx.Next()
	}
}
