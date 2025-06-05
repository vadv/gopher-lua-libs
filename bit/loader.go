package bit

import lua "github.com/yuin/gopher-lua"

// Preload adds bit to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local bit = require("bit")
func Preload(l *lua.LState) {
	l.PreloadModule("bit", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"band":   Bitwise(and),
	"bor":    Bitwise(or),
	"bxor":   Bitwise(xor),
	"lshift": Bitwise(ls),
	"rshift": Bitwise(rs),
	"bnot":   Not,
}
