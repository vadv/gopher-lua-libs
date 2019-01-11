// +build !windows
// +build sqlite

package db

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type luaSQLite struct {
	sync.Mutex
	db *sql.DB
}

func init() {
	RegisterDriver(`sqlite3`, &luaSQLite{})
}

func (sqlite *luaSQLite) constructor(connString string) (luaDB, error) {
	db, err := sql.Open(`sqlite3`, connString)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
	result := &luaSQLite{}
	result.db = db
	return result, nil
}

func (sqlite *luaSQLite) getDB() *sql.DB {
	sqlite.Lock()
	defer sqlite.Unlock()
	return sqlite.db
}
