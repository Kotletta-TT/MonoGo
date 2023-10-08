package http

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/gin-gonic/gin"
	"time"
)

func RequestResponseLogging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		duration := time.Since(startTime)
		WriteSize := ctx.Writer.Size()
		if WriteSize < 0 {
			WriteSize = 0
		}
		logger.Logger.Infof("%s %s %d %s %s %d bytes", ctx.ClientIP(), ctx.Request.Method, ctx.Writer.Status(), ctx.Request.URL.Path, duration, WriteSize)
	}
}
