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
	Compress       string `env:"COMPRESS"`
	BatchSupport   bool   `env:"BATCH_SUPPORT"`
}

func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "Address:port server")
	flag.IntVar(&config.ReportInterval, "r", 10, "Frequency to send server in sec")
	flag.IntVar(&config.PollInterval, "p", 2, "Frequency collect metrics in sec")
	flag.StringVar(&config.SendType, "s", "json", "Send type")
	flag.StringVar(&config.Compress, "compress", "gzip", "Compress send JSON-data")
	flag.BoolVar(&config.BatchSupport, "batch", false, "Use batch mode")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
