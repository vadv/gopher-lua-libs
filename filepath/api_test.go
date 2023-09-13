package filepath

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/inspect"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		Preload,
		inspect.Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
