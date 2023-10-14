package http

import (
	"compress/gzip"
	"io"
	"net/http"

	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/gin-gonic/gin"
)

const GZIP = "gzip"

type gzipResponseWriter struct {
	gin.ResponseWriter
	writer io.WriteCloser
}

func (w gzipResponseWriter) Write(data []byte) (int, error) {
	return w.writer.Write(data)
}

func (w gzipResponseWriter) Close() error {
	return w.writer.Close()
}

func CompressMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.GetHeader("Accept-Encoding") == GZIP {
			ctx.Writer.Header().Set("Content-Encoding", GZIP)
			gzipWriter := gzip.NewWriter(ctx.Writer)
			defer func() {
				if err := gzipWriter.Close(); err != nil {
					logger.Errorf("gzip writer close error: %s", err.Error())
				}
			}()
			gzippedResponseWriter := &gzipResponseWriter{ctx.Writer, gzipWriter}
			ctx.Writer = gzippedResponseWriter
		}
		if ctx.GetHeader("Content-Encoding") == GZIP {
			gzipReader, err := gzip.NewReader(ctx.Request.Body)
			if err != nil {
				logger.Errorf("gzip error: %s", err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			defer func() {
				if err := gzipReader.Close(); err != nil {
					logger.Errorf("gzip reader close error: %s", err.Error())
				}
			}()
			ctx.Request.Body = gzipReader
		}
		ctx.Next()
	}
}
