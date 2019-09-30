package cloudwatch

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds cloudwatch to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local cloudwatch = require("cloudwatch")
func Preload(L *lua.LState) {
	L.PreloadModule("cloudwatch", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	clwUd := L.NewTypeMetatable(`clw_ud`)
	L.SetGlobal(`clw_ud`, clwUd)
	L.SetField(clwUd, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"download":        Download,
		"get_metric_data": GetMetricData,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"new": New,
}
