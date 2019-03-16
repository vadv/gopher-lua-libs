// +build !windows
// +build sqlite

package db

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type luaSQLite struct {
	config *dbConfig
	sync.Mutex
	db *sql.DB
}

func init() {
	RegisterDriver(`sqlite3`, &luaSQLite{})
}

var (
	sharedSqlite     = make(map[string]*luaSQLite, 0)
	sharedSqliteLock = &sync.Mutex{}
)

func (sqlite *luaSQLite) constructor(config *dbConfig) (luaDB, error) {

	sharedSqliteLock.Lock()
	defer sharedSqliteLock.Unlock()

	if config.sharedMode {
		result, ok := sharedSqlite[config.connString]
		if ok {
			return result, nil
		}
	}

	db, err := sql.Open(`sqlite3`, config.connString)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(config.maxOpenConns)
	db.SetMaxOpenConns(config.maxOpenConns)
	result := &luaSQLite{config: config}
	result.db = db

	if config.sharedMode {
		sharedSqlite[config.connString] = result
	}

	return result, nil
}

func (sqlite *luaSQLite) getDB() *sql.DB {
	sqlite.Lock()
	defer sqlite.Unlock()
	return sqlite.db
}

func (sqlite *luaSQLite) getTXOptions() *sql.TxOptions {
	return &sql.TxOptions{ReadOnly: sqlite.config.readOnly}
}

func (sqlite *luaSQLite) closeDB() error {
	sqlite.Lock()
	defer sqlite.Unlock()
	err := sqlite.db.Close()
	if err != nil {
		return err
	}
	if sqlite.config.sharedMode {
		sharedSqliteLock.Lock()
		delete(sharedSqlite, sqlite.config.connString)
		sharedSqliteLock.Unlock()
	}
	return nil
}
