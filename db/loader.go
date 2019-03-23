package db

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds db to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local db = require("db")
func Preload(L *lua.LState) {
	L.PreloadModule("db", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	db_ud := L.NewTypeMetatable(`db_ud`)
	L.SetGlobal(`db_ud`, db_ud)
	L.SetField(db_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"query":   Query,
		"exec":    Exec,
		"stmt":    Stmt,
		"command": Command,
		"close":   Close,
	}))

	stmt_ud := L.NewTypeMetatable(`stmt_ud`)
	L.SetGlobal(`stmt_ud`, stmt_ud)
	L.SetField(stmt_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"query": StmtQuery,
		"exec":  StmtExec,
		"close": StmtClose,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"open": Open,
}
