package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

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
	go app.Run(cfg)
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
