package optparse

//go:generate go run internal/include_all_lua.go

import (
	"encoding/base64"

	lua "github.com/yuin/gopher-lua"
)

// Preload adds optparse to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local inspect = require("optparse")
func Preload(L *lua.LState) {
	L.PreloadModule("optparse", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	code, err := base64.StdEncoding.DecodeString(lua_optparse)
	if err != nil {
		panic(err.Error())
	}
	if err := L.DoString(string(code)); err != nil {
		L.RaiseError("load library 'optparse' error: %s", err.Error())
	}
	return 1
}
