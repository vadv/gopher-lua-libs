package db

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
)

type luaPG struct {
	sync.Mutex
	db *sql.DB
}

func init() {
	RegisterDriver(`postgres`, &luaPG{})
}

func (pg *luaPG) constructor(connString string) (luaDB, error) {
	db, err := sql.Open(`postgres`, connString)
	if err != nil {
		return nil, err
	}
	result := &luaPG{}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	result.db = db
	return result, nil
}

func (pg *luaPG) getDB() *sql.DB {
	pg.Lock()
	defer pg.Unlock()
	return pg.db
}
