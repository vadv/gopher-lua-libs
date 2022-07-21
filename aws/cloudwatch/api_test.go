package cloudwatch

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"

	"github.com/vadv/gopher-lua-libs/inspect"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		inspect.Preload,
		Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
