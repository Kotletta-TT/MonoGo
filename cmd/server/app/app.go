package app

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"log"
)

func Run() {
	parseFlags()
	memRepo := repository.New()
	ginRouter := http.NewRouter(memRepo)
	log.Fatal(ginRouter.Run(flagRunAddr))
}
