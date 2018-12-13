// build +purego

package db

import (
	"database/sql"
	"sync"

	_ "modernc.org/ql/driver"
)

type luaQL struct {
	sync.Mutex
	db *sql.DB
}

func init() {
	RegisterDriver(`ql`, &luaQL{})
	RegisterDriver(`ql2`, &luaQL{})
	RegisterDriver(`ql-mem`, &luaQL{})
}

func (ql *luaQL) constructor(connString string) (luaDB, error) {
	db, err := sql.Open(`ql2`, connString)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	result := &luaQL{}
	result.db = db
	return result, nil
}

func (ql *luaQL) getDB() *sql.DB {
	ql.Lock()
	defer ql.Unlock()
	return ql.db
}
