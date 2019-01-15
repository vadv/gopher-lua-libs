// +build !windows
// +build sqlite

package db

import (
	"testing"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	state := lua.NewState()

	Preload(state)
	time.Preload(state)
	inspect.Preload(state)

	if err := state.DoFile("./test/test_api.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
