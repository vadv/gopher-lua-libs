package plugin

import (
	"testing"

	ioutil "github.com/vadv/gopher-lua-libs/ioutil"
	time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	ioutil.Preload(state)
	if err := state.DoFile("./test/test_api.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
