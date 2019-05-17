// Package db implements golang package db functionality for lua.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

const (
	// max open connections
	MaxOpenConns = 1
)

type luaDB interface {
	constructor(*dbConfig) (luaDB, error)
	getDB() *sql.DB
	closeDB() error
	getTXOptions() *sql.TxOptions
}

type dbConfig struct {
	connString   string
	sharedMode   bool
	maxOpenConns int
	readOnly     bool
}

var (
	knownDrivers     = make(map[string]luaDB, 0)
	knownDriversLock = &sync.Mutex{}
	sharedDB         = make(map[string]luaDB, 0)
	sharedDBLock     = &sync.Mutex{}
)

// RegisterDriver register sql driver
func RegisterDriver(driver string, i luaDB) {
	knownDriversLock.Lock()
	defer knownDriversLock.Unlock()

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

// Open lua db.open(driver, connection_string, config) returns (db_ud, err)
// config table:
//   {
//     shared=false,
//     max_connections=X,
//     read_only=false
//   }
func Open(L *lua.LState) int {
	knownDriversLock.Lock()
	defer knownDriversLock.Unlock()

	driver := L.CheckString(1)
	connString := L.CheckString(2)
	db, ok := knownDrivers[driver]
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("unknown driver: %s", driver)))
		return 2
	}

	// parse config
	config := &dbConfig{connString: connString, maxOpenConns: MaxOpenConns}
	if L.GetTop() > 2 {
		configLua := L.CheckTable(3)
		configLua.ForEach(func(k lua.LValue, v lua.LValue) {
			if k.String() == `shared` {
				if val, ok := v.(lua.LBool); ok {
					config.sharedMode = bool(val)
				} else {
					L.ArgError(3, "shared must be bool")
				}
			}
			if k.String() == `max_connections` {
				if val, ok := v.(lua.LNumber); ok {
					config.maxOpenConns = int(val)
				} else {
					L.ArgError(3, "max_connections must be number")
				}
			}
			if k.String() == `read_only` {
				if val, ok := v.(lua.LBool); ok {
					config.readOnly = bool(val)
				} else {
					L.ArgError(3, "read_only must be bool")
				}
			}
		})
	}

	dbIface, err := db.constructor(config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	ud := L.NewUserData()
	ud.Value = dbIface
	L.SetMetatable(ud, L.GetTypeMetatable(`db_ud`))
	L.Push(ud)
	return 1
}

// Query lua db_ud:query(query) returns ({rows = {}, columns = {}}, err)
func Query(L *lua.LState) int {
	dbInterface := checkDB(L, 1)
	query := L.CheckString(2)
	sqlDB := dbInterface.getDB()
	opts := dbInterface.getTXOptions()
	tx, err := sqlDB.BeginTx(context.Background(), opts)
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

// Exec lua db_ud:exec(query) returns ({rows_affected=number, last_insert_id=number}, err)
func Exec(L *lua.LState) int {
	dbInterface := checkDB(L, 1)
	query := L.CheckString(2)
	sqlDB := dbInterface.getDB()
	opts := dbInterface.getTXOptions()
	tx, err := sqlDB.BeginTx(context.Background(), opts)
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

// Command lua db_ud:command(query) returns ({rows = {}, columns = {}}, err)
func Command(L *lua.LState) int {
	dbInterface := checkDB(L, 1)
	query := L.CheckString(2)
	sqlDB := dbInterface.getDB()
	sqlRows, err := sqlDB.Query(query)
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

// Close lua db_ud:close() returns err
func Close(L *lua.LState) int {
	dbIface := checkDB(L, 1)
	if err := dbIface.closeDB(); err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
