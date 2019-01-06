package chef

import (
	"testing"

	http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	state := lua.NewState()
	Preload(state)
	http.Preload(state)
	inspect.Preload(state)
	if err := state.DoFile("./test/test_api.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
