package main

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/app"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	log "github.com/Kotletta-TT/MonoGo/internal/server/logger"
)

func main() {
	cfg := config.NewConfig()
	log.Init(cfg)
	app.Run(cfg)
}
