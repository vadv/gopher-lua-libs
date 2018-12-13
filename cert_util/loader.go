package cert_util

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds cert_util to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local cert_util = require("cert_util")
func Preload(L *lua.LState) {
	L.PreloadModule("cert_util", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"not_after": NotAfter,
}
