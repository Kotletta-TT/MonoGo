package main

import (
	"log"

	"github.com/Kotletta-TT/MonoGo/cmd/agent/app"
	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
)

func main() {
	log.Println("Start Agent MonoGo")
	cnf := config.NewConfig()
	app.Run(cnf)
}
