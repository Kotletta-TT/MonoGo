package http

import (
	"compress/gzip"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
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
					logger.Logger.Infof("gzip writer close error: %s", err.Error())
				}
			}()
			gzippedResponseWriter := &gzipResponseWriter{ctx.Writer, gzipWriter}
			ctx.Writer = gzippedResponseWriter
		}
		if ctx.GetHeader("Content-Encoding") == GZIP {
			gzipReader, err := gzip.NewReader(ctx.Request.Body)
			defer func() {
				if err := gzipReader.Close(); err != nil {
					logger.Logger.Infof("gzip reader close error: %s", err.Error())
				}
			}()
			if err != nil {
				logger.Logger.Infof("gzip error: %s", err.Error())
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ctx.Request.Body = gzipReader
		}
		ctx.Next()
	}
}
