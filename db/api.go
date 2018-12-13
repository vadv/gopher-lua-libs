// Package db implements golang package db functionality for lua.
package db

import (
	"database/sql"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

const (
	// max idle connections
	MaxIdleConns = 1
	// max open connections
	MaxOpenConns = 1
)

type luaDB interface {
	constructor(string) (luaDB, error)
	getDB() *sql.DB
}

var knownDrivers = make(map[string]luaDB, 0)

// RegisterDriver(): register sql driver
func RegisterDriver(driver string, i luaDB) {
	knownDrivers[driver] = i
}

func checkDB(L *lua.LState, n int) luaDB {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(luaDB); ok {
		return v
	}
	L.ArgError(n, "database expected")
	return nil
}

// Open(): lua db:open(driver, connection_string) returns (db_ud, err)
func Open(L *lua.LState) int {
	driver := L.CheckString(1)
	connString := L.CheckString(2)
	db, ok := knownDrivers[driver]
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unknown driver: %s", driver)))
		return 2
	}
	result, err := db.constructor(connString)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = result
	L.SetMetatable(ud, L.GetTypeMetatable(`db_ud`))
	L.Push(ud)
	return 1
}

// Query(): lua db_ud:query(query) returns ({rows = {}, columns = {}}, err)
func Query(L *lua.LState) int {
	dbInterface := checkDB(L, 1)
	query := L.CheckString(2)
	sqlDB := dbInterface.getDB()
	tx, err := sqlDB.Begin()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer tx.Commit()
	sqlRows, err := tx.Query(query)
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

// Exec(): lua db_ud:exec(query) returns ({rows_affected=number, last_insert_id=number}, err)
func Exec(L *lua.LState) int {
	dbInterface := checkDB(L, 1)
	query := L.CheckString(2)
	sqlDB := dbInterface.getDB()
	tx, err := sqlDB.Begin()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer tx.Commit()
	sqlResult, err := tx.Exec(query)
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

// Close(): lua db_ud:query(query) returns (table of tables, err)
func Close(L *lua.LState) int {
	dbInterface := checkDB(L, 1)
	sqlDB := dbInterface.getDB()
	if err := sqlDB.Close(); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
