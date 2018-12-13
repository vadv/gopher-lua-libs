package xmlpath

import (
	"io/ioutil"
	"testing"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	data, err := ioutil.ReadFile("./test/test_api.lua")
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	if err := state.DoString(string(data)); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
