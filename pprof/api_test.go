package pprof_test

import (
	"github.com/stretchr/testify/assert"
	lua_http "github.com/vadv/gopher-lua-libs/http"
	lua_pprof "github.com/vadv/gopher-lua-libs/pprof"
	"github.com/vadv/gopher-lua-libs/tests"
	lua_time "github.com/vadv/gopher-lua-libs/time"
	"testing"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		lua_pprof.Preload,
		lua_http.Preload,
		lua_time.Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
