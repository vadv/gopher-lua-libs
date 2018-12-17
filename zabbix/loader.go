package zabbix

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds zabbix to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local zabbix = require("zabbix")
func Preload(L *lua.LState) {
	L.PreloadModule("zabbix", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	zabbix_bot_ud := L.NewTypeMetatable(`zabbix_bot_ud`)
	L.SetGlobal(`zabbix_bot_ud`, zabbix_bot_ud)
	L.SetField(zabbix_bot_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"login":      Login,
		"logout":     Logout,
		"request":    Request,
		"save_graph": SaveGraph,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"bot": NewBot,
	"new": NewBot,
}
