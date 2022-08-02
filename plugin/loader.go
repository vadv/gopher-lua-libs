package plugin

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds plugin to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local plugin = require("plugin")
func Preload(L *lua.LState) {
	L.PreloadModule("plugin", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	pluginUd := L.NewTypeMetatable(`plugin_ud`)
	L.SetGlobal(`plugin_ud`, pluginUd)
	L.SetField(pluginUd, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"run":          Run,
		"error":        Error,
		"stop":         Stop,
		"wait":         Wait,
		"is_running":   IsRunning,
		"done_channel": DoneChannel,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"do_string":              DoString,
	"do_file":                DoFile,
	"do_string_with_payload": DoStringWithPayload,
	"do_file_with_payload":   DoFileWithPayload,
}
