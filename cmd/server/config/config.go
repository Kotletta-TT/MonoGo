package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	RunServerAddr string `env:"ADDRESS"`
}

func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.RunServerAddr, "a", "localhost:8080", "Address:port server")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
