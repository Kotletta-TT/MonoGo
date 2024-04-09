// Package app implements the main process of the application.
package app

import (
	"context"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	"github.com/Kotletta-TT/MonoGo/internal/server/controller/http"
	log "github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage"
	"golang.org/x/sync/errgroup"
)

// Run executes the Go function.
//
// It takes a pointer to a `config.Config` struct as a parameter.
// It does not return anything.
func Run(ctx context.Context, cfg *config.Config) {
	repo := storage.GetRepo(ctx, cfg)
	defer repo.Close()
	srv, err := http.NewServer(repo, cfg)
	if err != nil {
		log.Error(err)
		return
	}
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if cfg.SSL {
			return srv.ListenAndServeTLS("", "")
		}
		return srv.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()
		return srv.Shutdown(context.Background())
	})
	log.Info(g.Wait())
}
