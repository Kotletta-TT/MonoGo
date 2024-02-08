// Package config implements a configuration object for the server.
package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	RunServerAddr   string `env:"ADDRESS" json:"address"`
	LogLevel        string `env:"LOG_LEVEL" json:"log_level"`
	LogPath         string `env:"LOG_PATH" json:"log_path"`
	LogFile         bool   `env:"LOG_FILE" json:"log_file"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DatabaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	HashKey         string `env:"KEY" json:"hash_key"`
	SSL             bool   `env:"SSL" json:"ssl"`
	CertPath        string `env:"CERT_PATH" json:"cert_path"`
	KeyPath         string `env:"KEY_PATH" json:"key_path"`
	CaPath          string `env:"CA_PATH" json:"ca_path"`
}

// NewConfig initializes a new Config object with default values and parses command line arguments and environment variables to override the defaults.
//
// Returns a pointer to the newly created Config object.
func NewConfig() *Config {
	config := Config{}
	var configPath string
	flag.StringVar(&config.RunServerAddr, "a", "localhost:8080", "Address:port server")
	flag.StringVar(&config.LogLevel, "l", "INFO", "Log level")
	flag.StringVar(&config.LogPath, "p", "/var/log/monogo.log", "Log path")
	flag.BoolVar(&config.LogFile, "log-file", false, "Log file")
	flag.IntVar(&config.StoreInterval, "i", 300, "Frequency to store server in sec")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "File storage path")
	flag.BoolVar(&config.Restore, "r", true, "Restore from file")
	flag.StringVar(&config.DatabaseDSN, "d", "", "DB URL example: postgres://username:password@localhost:5432/database_name")
	flag.StringVar(&config.HashKey, "k", "", "Hash key for signing data")
	flag.BoolVar(&config.SSL, "s", false, "Use SSL")
	flag.StringVar(&config.CertPath, "cert", "server.crt", "Path to certificate")
	flag.StringVar(&config.KeyPath, "key", "server.key", "Path to key")
	flag.StringVar(&config.CaPath, "ca", "root.pem", "Path to certificate authority")
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
