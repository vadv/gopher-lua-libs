package prometheus_client

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds prometheus to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local prometheus = require("prometheus")
func Preload(L *lua.LState) {
	L.PreloadModule("prometheus", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	prometheusClient := L.NewTypeMetatable(`prometheus_client_ud`)
	L.SetGlobal(`prometheus_client_ud`, prometheusClient)
	L.SetField(prometheusClient, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"start": Start,
		"stop":  Stop,
	}))

	prometheusMetricUd := L.NewTypeMetatable(`prometheus_client_metric_ud`)
	L.SetGlobal(`prometheus_client_metric_ud`, prometheusMetricUd)
	L.SetField(prometheusMetricUd, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set": Set,
		"add": Add,
		"inc": Inc,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"register": Register,
	"counter":  Counter,
	"gauge":    Gauge,
}
