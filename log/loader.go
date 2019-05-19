package log

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds log package to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local log = require("log")
func Preload(L *lua.LState) {
	L.PreloadModule("log", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	logger := L.NewTypeMetatable(`logger`)
	L.SetGlobal(`logger`, logger)
	L.SetField(logger, `__index`, L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		`debug`: Debug,
		`info`:  Info,
		`warn`:  Warn,
		`error`: Error,
		`fatal`: Fatal,
	}))
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	`logger`: NewLogger,
}
