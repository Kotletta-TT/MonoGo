package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGzipCompress_Empty(t *testing.T) {
	zipData, err := GzipCompress([]byte(""))
	testData := []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 1, 0, 0, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0}
	assert.NoError(t, err)
	assert.Equal(t, testData, zipData)
}

func TestGzipCompress_SomeData(t *testing.T) {
	zipData, err := GzipCompress([]byte("some data"))
	testData := []byte{31, 139, 8, 0, 0, 0, 0, 0, 2, 255, 42, 206, 207, 77, 85, 72, 73, 44, 73, 4, 4, 0, 0, 255, 255, 30, 233, 194, 217, 9, 0, 0, 0}
	assert.NoError(t, err)
	assert.Equal(t, testData, zipData)
}
