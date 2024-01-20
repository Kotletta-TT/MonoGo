package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Foo() {}

func TestGetFunctionName(t *testing.T) {
	name := GetFunctionName(Foo)
	assert.Equal(t, "github.com/Kotletta-TT/MonoGo/internal/agent/utils.Foo", name)
}
