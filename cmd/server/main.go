package main

import (
	"github.com/Kotletta-TT/MonoGo/cmd/server/app"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
)

func main() {
	cfg := config.NewConfig()
	app.Run(cfg)
}
