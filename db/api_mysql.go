package db

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type luaMySQL struct {
	sync.Mutex
	db *sql.DB
}

func init() {
	RegisterDriver(`mysql`, &luaMySQL{})
}

func (mysql *luaMySQL) constructor(connString string) (luaDB, error) {
	db, err := sql.Open(`mysql`, connString)
	if err != nil {
		return nil, err
	}
	result := &luaMySQL{}
	db.SetMaxIdleConns(MaxIdleConns)
	db.SetMaxOpenConns(MaxOpenConns)
	result.db = db
	return result, nil
}

func (mysql *luaMySQL) getDB() *sql.DB {
	mysql.Lock()
	defer mysql.Unlock()
	return mysql.db
}
