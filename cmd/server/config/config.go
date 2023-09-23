package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	RunServerAddr string `env:"ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
	LogPath       string `env:"LOG_PATH"`
}

func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.RunServerAddr, "a", "localhost:8080", "Address:port server")
	flag.StringVar(&config.LogLevel, "l", "INFO", "Log level")
	flag.StringVar(&config.LogPath, "p", "/var/log/monogo.log", "Log path")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
