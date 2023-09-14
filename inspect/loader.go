package inspect

import (
	_ "embed"
	lua "github.com/yuin/gopher-lua"
)

//go:embed inspect.lua
var lua_inspect string

// Preload adds inspect to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local inspect = require("inspect")
func Preload(L *lua.LState) {
	L.PreloadModule("inspect", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	if err := L.DoString(lua_inspect); err != nil {
		L.RaiseError("load library 'inspect' error: %s", err.Error())
	}
	return 1
}
