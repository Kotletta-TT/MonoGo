// Package config implements a configuration object for the agent.
package config

import (
	"flag"
	"runtime"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	ServerHost     string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	SendType       string `env:"SEND_TYPE"`
	Compress       string `env:"COMPRESS"`
	BatchSupport   bool   `env:"BATCH_SUPPORT"`
	HashKey        string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	SSL            bool   `env:"CRYPTO_KEY"`
	CertPath       string `env:"CERT_PATH"`
	KeyPath        string `env:"KEY_PATH"`
	CaPath         string `env:"CA_PATH"`
}

// NewConfig creates a new Config object and initializes its fields
// based on command-line flags and environment variables.
//
// No parameters.
// Returns a pointer to the newly created Config object.
func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.ServerHost, "a", "localhost:8080", "Address:port server")
	flag.IntVar(&config.ReportInterval, "r", 10, "Frequency to send server in sec")
	flag.IntVar(&config.PollInterval, "p", 2, "Frequency collect metrics in sec")
	flag.StringVar(&config.SendType, "t", "json", "Send type")
	flag.StringVar(&config.Compress, "compress", "gzip", "Compress send JSON-data")
	flag.BoolVar(&config.BatchSupport, "batch", false, "Use batch mode")
	flag.StringVar(&config.HashKey, "k", "", "Hash key for signing data")
	flag.IntVar(&config.RateLimit, "l", runtime.NumCPU(), "Rate limit")
	flag.BoolVar(&config.SSL, "s", false, "Use SSL")
	flag.StringVar(&config.CertPath, "cert", "agent.crt", "Certificate path")
	flag.StringVar(&config.KeyPath, "key", "agent.key", "Key path")
	flag.StringVar(&config.CaPath, "ca", "root.pem", "CA path IMPORTANT! use pem format")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
