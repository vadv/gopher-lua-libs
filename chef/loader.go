package chef

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds chef to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local chef = require("chef")
func Preload(L *lua.LState) {
	L.PreloadModule("chef", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	chef_client_ud := L.NewTypeMetatable(`chef_client_ud`)
	L.SetGlobal(`chef_client_ud`, chef_client_ud)
	L.SetField(chef_client_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"request": Request,
		"search":  Search,
	}))
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"client": NewClient,
}
