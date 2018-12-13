package tcp

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds tcp to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local tcp = require("tcp")
func Preload(L *lua.LState) {
	L.PreloadModule("tcp", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	tcp_client_ud := L.NewTypeMetatable(`tcp_client_ud`)
	L.SetGlobal(`tcp_client_ud`, tcp_client_ud)
	L.SetField(tcp_client_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"write": Write,
		"close": Close,
		"read":  Read,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"open": Open,
}
