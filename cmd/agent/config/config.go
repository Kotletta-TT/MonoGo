package config

import (
	"flag"
	"github.com/caarlos0/env/v9"
)

type Config struct {
	ServerHost     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	SendType       string `env:"SEND_TYPE"`
}

func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "Address:port server")
	flag.IntVar(&config.ReportInterval, "r", 10, "Frequency to send server in sec")
	flag.IntVar(&config.PollInterval, "p", 2, "Frequency collect metrics in sec")
	flag.StringVar(&config.SendType, "s", "json", "Send type")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
