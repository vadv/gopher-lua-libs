package pprof_test

import (
	"testing"

	lua_http "github.com/vadv/gopher-lua-libs/http"
	lua_pprof "github.com/vadv/gopher-lua-libs/pprof"
	lua_time "github.com/vadv/gopher-lua-libs/time"

	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	state := lua.NewState()
	lua_pprof.Preload(state)
	lua_http.Preload(state)
	lua_time.Preload(state)
	if err := state.DoFile("./test/test_api.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
