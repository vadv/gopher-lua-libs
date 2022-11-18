package log

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/filepath"
	"github.com/vadv/gopher-lua-libs/strings"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"

	ioutil "github.com/vadv/gopher-lua-libs/ioutil"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		ioutil.Preload,
		Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}

func TestLogLevelApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		ioutil.Preload,
		filepath.Preload,
		strings.Preload,
		Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_loglevel.lua"))
}
