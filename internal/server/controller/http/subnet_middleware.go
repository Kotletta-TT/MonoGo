package http

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/netip"
)

func TrustedSubnetMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if cfg.TrustedSubnet == "" {
			ctx.Next()
			return
		}
		network, err := netip.ParsePrefix(cfg.TrustedSubnet)
		if err != nil {
			logger.Error(err.Error())
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		addr, err := netip.ParseAddr(ctx.GetHeader("X-Real-IP"))
		if err != nil {
			logger.Error(err.Error())
			ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if !network.Contains(addr) {
			err = fmt.Errorf("address %s not allowed", addr.String())
			logger.Error(err)
			ctx.JSON(http.StatusForbidden, gin.H{"error": err})
			return
		}
		ctx.Next()
	}
}
