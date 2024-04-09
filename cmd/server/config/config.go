// Package config implements a configuration object for the server.
package config

import (
	"flag"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	RunServerAddr   string `env:"ADDRESS"`
	LogLevel        string `env:"LOG_LEVEL"`
	LogPath         string `env:"LOG_PATH"`
	LogFile         bool   `env:"LOG_FILE"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	HashKey         string `env:"KEY"`
}

// NewConfig initializes a new Config object with default values and parses command line arguments and environment variables to override the defaults.
//
// Returns a pointer to the newly created Config object.
func NewConfig() *Config {
	config := Config{}
	flag.StringVar(&config.RunServerAddr, "a", "localhost:8080", "Address:port server")
	flag.StringVar(&config.LogLevel, "l", "INFO", "Log level")
	flag.StringVar(&config.LogPath, "p", "/var/log/monogo.log", "Log path")
	flag.BoolVar(&config.LogFile, "log-file", false, "Log file")
	flag.IntVar(&config.StoreInterval, "i", 300, "Frequency to store server in sec")
	flag.StringVar(&config.FileStoragePath, "f", "/tmp/metrics-db.json", "File storage path")
	flag.BoolVar(&config.Restore, "r", true, "Restore from file")
	flag.StringVar(&config.DatabaseDSN, "d", "", "DB URL example: postgres://username:password@localhost:5432/database_name")
	flag.StringVar(&config.HashKey, "k", "", "Hash key for signing data")
	flag.Parse()
	err := env.Parse(&config)
	if err != nil {
		panic(err)
	}
	return &config
}
