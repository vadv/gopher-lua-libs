package argparse

import (
	lua "github.com/yuin/gopher-lua"
)

// TODO(scr): move to embed once minimum supported go version is 1.16
//go:generate go run github.com/logrusorgru/textFileToGoConst@latest -in argparse.lua -o lua_const.go -c lua_argparse

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
