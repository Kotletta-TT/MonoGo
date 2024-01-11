package main

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/app"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	log "github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	cfg := config.NewConfig()
	log.Init(cfg)
	go app.Run(cfg)
	http.ListenAndServe("localhost:8080", nil)
}
