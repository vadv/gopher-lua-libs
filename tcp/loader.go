package tcp

import (
	lua "github.com/yuin/gopher-lua"
	"time"
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

	funcs := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"write": Write,
		"close": Close,
		"read":  Read,
	})
	L.SetFuncs(tcp_client_ud, map[string]lua.LGFunction{
		"__index": func(state *lua.LState) int {
			conn := checkLuaTCPClient(L, 1)
			k := L.CheckString(2)
			var duration time.Duration
			switch k {
			case "dialTimeout":
				duration = conn.dialTimeout
			case "writeTimeout":
				duration = conn.writeTimeout
			case "readTimeout":
				duration = conn.readTimeout
			case "closeTimeout":
				duration = conn.closeTimeout
			default:
				L.Push(L.GetField(funcs, k))
				return 1
			}
			L.Push(lua.LNumber(duration) / lua.LNumber(time.Second))
			return 1
		},
		"__newindex": func(state *lua.LState) int {
			conn := checkLuaTCPClient(L, 1)
			k := L.CheckString(2)
			var pDuration *time.Duration
			switch k {
			case "dialTimeout":
				pDuration = &conn.dialTimeout
			case "writeTimeout":
				pDuration = &conn.writeTimeout
			case "readTimeout":
				pDuration = &conn.readTimeout
			case "closeTimeout":
				pDuration = &conn.closeTimeout
			default:
				return 0
			}
			*pDuration = time.Duration(L.CheckNumber(3) * lua.LNumber(time.Second))
			return 0
		},
	})

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"open": Open,
}
