package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"log"
)

type Config struct {
	ServerHost     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POLL_INTERVAL"`
}

func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "Address:port server")
	flag.IntVar(&config.ReportInterval, "r", 10, "Frequency to send server in sec")
	flag.IntVar(&config.PoolInterval, "p", 2, "Frequency collect metrics in sec")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	log.Print(config)
	return &config
}
