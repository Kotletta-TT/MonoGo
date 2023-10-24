package app

import (
	"database/sql"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	"github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Run(cfg *config.Config) {
	memRepo := storage.New(cfg)
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		logger.Fatalf("Connect to database error: %s connect to: %s", err, cfg.DatabaseDSN)
	}
	go memRepo.Stash()
	ginRouter := http.NewRouter(memRepo, db)
	logger.Infof("Start server: http://%s/", cfg.RunServerAddr)
	logger.Fatal(ginRouter.Run(cfg.RunServerAddr))
}
