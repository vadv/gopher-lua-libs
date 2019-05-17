package db

import (
	"database/sql"

	lua "github.com/yuin/gopher-lua"
)

type luaStmt struct {
	*sql.Stmt
	d *sql.DB
}

// Stmt lua db_ud:stmt(query) returns (stmt_ud, err)
func Stmt(L *lua.LState) int {
	dbInterface := checkDB(L, 1)
	query := L.CheckString(2)
	sqlDB := dbInterface.getDB()
	s, err := sqlDB.Prepare(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = &luaStmt{Stmt: s, d: sqlDB}
	L.SetMetatable(ud, L.GetTypeMetatable(`stmt_ud`))
	L.Push(ud)
	return 1
}

func getSTMTArgs(L *lua.LState) []interface{} {
	args := make([]interface{}, 0)
	for i := 2; i <= L.GetTop(); i++ {
		any := L.CheckAny(i)
		switch any.Type() {
		case lua.LTNil:
			args = append(args, nil)
		default:
			args = append(args, L.CheckAny(i))
		}
	}
	return args
}

// StmtQuery lua stmt_ud:query(args) returns ({rows = {}, columns = {}}, err)
func StmtQuery(L *lua.LState) int {
	ud := L.CheckUserData(1)
	s, ok := ud.Value.(*luaStmt)
	if !ok {
		L.ArgError(1, "must be stmt_ud")
	}
	args := getSTMTArgs(L)
	sqlRows, err := s.Query(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer sqlRows.Close()
	rows, columns, err := parseRows(sqlRows, L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	result.RawSetString(`rows`, rows)
	result.RawSetString(`columns`, columns)
	L.Push(result)
	return 1
}

// StmtExec lua stmt_ud:exec(args) returns ({rows_affected=number, last_insert_id=number}, err)
func StmtExec(L *lua.LState) int {
	ud := L.CheckUserData(1)
	s, ok := ud.Value.(*luaStmt)
	if !ok {
		L.ArgError(1, "must be stmt_ud")
	}
	args := getSTMTArgs(L)
	sqlResult, err := s.Exec(args...)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	if id, err := sqlResult.LastInsertId(); err == nil {
		result.RawSetString(`last_insert_id`, lua.LNumber(id))
	}
	if aff, err := sqlResult.RowsAffected(); err == nil {
		result.RawSetString(`rows_affected`, lua.LNumber(aff))
	}
	L.Push(result)
	return 1
}

// StmtClose lua stmt_ud:close() returns err
func StmtClose(L *lua.LState) int {
	ud := L.CheckUserData(1)
	s, ok := ud.Value.(*luaStmt)
	if !ok {
		L.ArgError(1, "must be stmt_ud")
	}
	if err := s.Close(); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
