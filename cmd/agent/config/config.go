// Package config implements a configuration object for the agent.
package config

import (
	"encoding/json"
	"flag"
	"os"
	"runtime"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	ServerHost     string `env:"ADDRESS" json:"address"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	SendType       string `env:"SEND_TYPE" json:"send_type"`
	Compress       string `env:"COMPRESS" json:"compress"`
	BatchSupport   bool   `env:"BATCH_SUPPORT" json:"batch_support"`
	HashKey        string `env:"KEY" json:"hash_key"`
	RateLimit      int    `env:"RATE_LIMIT" json:"rate_limit"`
	SSL            bool   `env:"CRYPTO_KEY" json:"ssl"`
	CertPath       string `env:"CERT_PATH" json:"cert_path"`
	KeyPath        string `env:"KEY_PATH" json:"key_path"`
	CaPath         string `env:"CA_PATH" json:"ca_path"`
}

// NewConfig creates a new Config object and initializes its fields
// based on command-line flags and environment variables.
//
// No parameters.
// Returns a pointer to the newly created Config object.
func NewConfig() *Config {
	config := Config{}
	var configPath string
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
	flag.StringVar(&configPath, "c", "", "Config path")
	flag.Parse()
	if configPath != "" {
		config = GetJSONConfig(configPath)
		flag.Parse()
	}
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}

func GetJSONConfig(configPath string) Config {
	f, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(f)
	config := new(Config)
	err = decoder.Decode(config)
	if err != nil {
		panic(err)
	}
	return *config
}
