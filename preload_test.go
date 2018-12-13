package libs

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestPreload(t *testing.T) {
	state := lua.NewState()
	Preload(state)
	if err := state.DoFile("./preload.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
