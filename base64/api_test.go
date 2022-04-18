package base64

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"
)

func TestApi(t *testing.T) {
	assert.NotZero(t, tests.RunLuaTestFile(t, Preload, "./test/test_api.lua"))
}
