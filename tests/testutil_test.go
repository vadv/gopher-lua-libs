package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/strings"
	"testing"
)

func TestSuite(t *testing.T) {
	preload := SeveralPreloadFuncs(
		strings.Preload,
	)
	assert.NotZero(t, RunLuaTestFile(t, preload, "testdata/suite_test.lua"))
}
