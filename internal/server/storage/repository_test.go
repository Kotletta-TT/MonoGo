package storage

import (
	"context"
	"testing"

	"github.com/Kotletta-TT/MonoGo/cmd/server/config"
	log "github.com/Kotletta-TT/MonoGo/internal/server/logger"
	"github.com/Kotletta-TT/MonoGo/internal/server/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestGetRepo(t *testing.T) {
	log.Init(&config.Config{LogLevel: "debug"})
	repo := GetRepo(context.Background(), &config.Config{DatabaseDSN: ""})
	assert.IsType(t, &memory.MemRepository{}, repo)
}
