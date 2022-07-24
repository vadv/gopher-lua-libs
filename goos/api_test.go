package goos

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"

	runtime "github.com/vadv/gopher-lua-libs/runtime"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		runtime.Preload,
		Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
