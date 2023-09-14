package log

import (
	_ "embed"
	lua "github.com/yuin/gopher-lua"
)

//go:embed loglevel.lua
var lua_loglevel string

// Preload adds log to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//		local log = require("log")
//	 or for levelled logging
//		local log = require("loglevel")
func Preload(L *lua.LState) {
	L.PreloadModule("log", Loader)
	L.PreloadModule("loglevel", LoadLogLevel)
}

func LoadLogLevel(L *lua.LState) int {
	if err := L.DoString(lua_loglevel); err != nil {
		L.RaiseError("load library 'loglevel' error: %s", err.Error())
	}
	return 1
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	loggerUD := L.NewTypeMetatable(`logger_ud`)
	L.SetGlobal(`logger_ud`, loggerUD)
	L.SetField(loggerUD, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set_output": SetOutput,
		"set_prefix": SetPrefix,
		"set_flags":  SetFlags,
		"print":      Print,
		"printf":     Printf,
		"println":    Println,
		"close":      Close,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"new": New,
}
