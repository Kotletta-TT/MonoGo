package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"

	"github.com/Kotletta-TT/MonoGo/cmd/server/app"
	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	log "github.com/Kotletta-TT/MonoGo/internal/server/logger"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	printBuildInfo()
	cfg := config.NewConfig()
	log.Init(cfg)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()
	go app.Run(ctx, cfg)
	http.ListenAndServe("localhost:8080", nil)
}

func printBuildInfo() {
	fmt.Printf("Version: %s", prettyStringBuildInfo(buildVersion))
	fmt.Printf("Date: %s", prettyStringBuildInfo(buildDate))
	fmt.Printf("Commit: %s", prettyStringBuildInfo(buildCommit))
}

func prettyStringBuildInfo(src string) string {
	if src == "" {
		return "N/A"
	}
	return src
}
