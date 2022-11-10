package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/goos"
	"github.com/vadv/gopher-lua-libs/strings"
	"testing"
)

func TestSuite(t *testing.T) {
	preload := strings.Preload
	assert.NotZero(t, RunLuaTestFile(t, preload, "testdata/test_suite.lua"))
}

func TestApi(t *testing.T) {
	preload := goos.Preload
	assert.NotZero(t, RunLuaTestFile(t, preload, "testdata/test_api.lua"))
}
