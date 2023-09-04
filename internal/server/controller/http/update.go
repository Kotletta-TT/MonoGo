package http

import (
	"errors"
	"github.com/Kotletta-TT/MonoGo/internal/server/infrastructure/repository"
	"github.com/Kotletta-TT/MonoGo/internal/server/usecase"
	"log"
	"net/http"
)

var incorrectTypeMetrics usecase.IncorrectTypeMetrics
var incorrectValueMetrics usecase.IncorrectValueMetrics
var noNameMetric usecase.NoNameMetric

type UpdateHandler struct {
	repo repository.Repository
}

func NewUpdateHandler(repo repository.Repository) UpdateHandler {
	return UpdateHandler{
		repo: repo,
	}
}

func (uh UpdateHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := usecase.Parse(req)
	if err != nil {
		switch {
		case errors.As(err, &incorrectTypeMetrics) || errors.As(err, &incorrectValueMetrics):
			res.WriteHeader(http.StatusBadRequest)
		case errors.As(err, &noNameMetric):
			res.WriteHeader(http.StatusNotFound)
		default:
			res.WriteHeader(http.StatusNotFound)
		}
		log.Println(err)
		return
	}

	res.WriteHeader(http.StatusOK)
}
