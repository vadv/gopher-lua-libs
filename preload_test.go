package libs

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"
)

func TestPreload(t *testing.T) {
	assert.NotZero(t, tests.RunLuaTestFile(t, Preload, "./preload.lua"))
}
