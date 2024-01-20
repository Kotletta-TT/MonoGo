// Package http implements some utils
package http

import (
	"fmt"
	"time"

	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/gin-gonic/gin"
)

// RequestResponseLogging returns a gin.HandlerFunc that logs the request and response details.
//
// It takes in a gin.Context as a parameter.
// There is no return value.
func RequestResponseLogging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		for header, value := range ctx.Request.Header {
			fmt.Printf("%s: %s\n", header, value)
		}
		ctx.Next()
		duration := time.Since(startTime)
		WriteSize := ctx.Writer.Size()
		if WriteSize < 0 {
			WriteSize = 0
		}
		logger.Infof("%s %s %d %s %s %d bytes", ctx.ClientIP(), ctx.Request.Method, ctx.Writer.Status(), ctx.Request.URL.Path, duration, WriteSize)
	}
}
