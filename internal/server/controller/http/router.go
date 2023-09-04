package http

import (
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"net/http"
)

func NewRouter(repo repository.Repository) *http.ServeMux {
	mux := http.NewServeMux()
	updateHandler := NewUpdateHandler(repo)
	mux.Handle("/update/", updateHandler)
	return mux
}
