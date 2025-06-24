package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/hex"
	"github.com/vadv/gopher-lua-libs/tests"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(Preload, hex.Preload)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
