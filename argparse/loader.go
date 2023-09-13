package argparse

import (
	_ "embed"
	lua "github.com/yuin/gopher-lua"
)

//go:embed argparse.lua
var lua_argparse string

// Preload adds inspect to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local inspect = require("inspect")
func Preload(L *lua.LState) {
	L.PreloadModule("argparse", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	if err := L.DoString(lua_argparse); err != nil {
		L.RaiseError("load library 'argparse' error: %s", err.Error())
	}
	return 1
}
