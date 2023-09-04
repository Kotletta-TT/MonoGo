package app

import (
	router "github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"log"
	"net/http"
)

func Run() {
	memRepo := repository.NewMemRepo()
	httpRouter := router.NewRouter(memRepo)
	log.Fatal(http.ListenAndServe(":8080", httpRouter))
}
