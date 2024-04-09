package main

import (
	"fmt"
	"github.com/Kotletta-TT/MonoGo/cmd/agent/app"
	"github.com/Kotletta-TT/MonoGo/cmd/agent/config"
	"log"
)

var buildVersion string
var buildDate string
var buildCommit string

func main() {
	printBuildInfo()
	log.Println("Start Agent MonoGo")
	cnf := config.NewConfig()
	app.Run(cnf)
}

func printBuildInfo() {
	fmt.Printf("Version: %s\n", prettyStringBuildInfo(buildVersion))
	fmt.Printf("Date: %s\n", prettyStringBuildInfo(buildDate))
	fmt.Printf("Commit: %s\n", prettyStringBuildInfo(buildCommit))
}

func prettyStringBuildInfo(src string) string {
	if src == "" {
		return "N/A"
	}
	return src
}
