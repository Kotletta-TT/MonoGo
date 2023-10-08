package app

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
)

func Run(cfg *config.Config) {
	memRepo := storage.New()
	ginRouter := http.NewRouter(memRepo)
	logger.Logger.Infof("Start server: http://%s/", cfg.RunServerAddr)
	logger.Logger.Fatal(ginRouter.Run(cfg.RunServerAddr))
}
