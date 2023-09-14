package app

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"log"
)

func Run(cfg *config.Config) {
	memRepo := storage.New()
	ginRouter := http.NewRouter(memRepo)
	log.Fatal(ginRouter.Run(cfg.RunServerAddr))
}
