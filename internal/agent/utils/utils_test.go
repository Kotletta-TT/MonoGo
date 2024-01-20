package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Foo() {}

func TestGetFunctionName(t *testing.T) {
	name := GetFunctionName(Foo)
	assert.Equal(t, "github.com/Kotletta-TT/MonoGo/internal/agent/utils.Foo", name)
}
