package tac

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds tac to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local tac = require("tac")
func Preload(L *lua.LState) {
	L.PreloadModule("tac", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	tac_ud := L.NewTypeMetatable(`tac_ud`)
	L.SetGlobal(`tac_ud`, tac_ud)
	L.SetField(tac_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"line":  Line,
		"close": Close,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"open": Open,
}
