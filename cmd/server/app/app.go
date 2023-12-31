package app

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
)

func Run(cfg *config.Config) {
	repo := storage.GetRepo(cfg)
	defer repo.Close()
	ginRouter := http.NewRouter(repo, cfg)
	logger.Infof("Start server: http://%s/", cfg.RunServerAddr)
	logger.Fatal(ginRouter.Run(cfg.RunServerAddr))
}
