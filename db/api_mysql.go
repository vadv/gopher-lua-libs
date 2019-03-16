package db

import (
	"database/sql"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type luaMySQL struct {
	config *dbConfig
	sync.Mutex
	db *sql.DB
}

var (
	sharedMySQL     = make(map[string]*luaMySQL, 0)
	sharedMySQLLock = &sync.Mutex{}
)

func init() {
	RegisterDriver(`mysql`, &luaMySQL{})
}

func (mysql *luaMySQL) constructor(config *dbConfig) (luaDB, error) {

	sharedMySQLLock.Lock()
	defer sharedMySQLLock.Unlock()

	if config.sharedMode {
		result, ok := sharedMySQL[config.connString]
		if ok {
			return result, nil
		}
	}

	db, err := sql.Open(`mysql`, config.connString)
	if err != nil {
		return nil, err
	}
	result := &luaMySQL{config: config}
	db.SetMaxIdleConns(config.maxOpenConns)
	db.SetMaxOpenConns(config.maxOpenConns)
	result.db = db

	if config.sharedMode {
		sharedMySQL[config.connString] = result
	}

	return result, nil
}

func (mysql *luaMySQL) getDB() *sql.DB {
	mysql.Lock()
	defer mysql.Unlock()
	return mysql.db
}

func (mysql *luaMySQL) getTXOptions() *sql.TxOptions {
	return &sql.TxOptions{ReadOnly: mysql.config.readOnly}
}

func (mysql *luaMySQL) closeDB() error {
	mysql.Lock()
	defer mysql.Unlock()
	err := mysql.db.Close()
	if err != nil {
		return err
	}
	if mysql.config.sharedMode {
		sharedMySQLLock.Lock()
		delete(sharedMySQL, mysql.config.connString)
		sharedMySQLLock.Unlock()
	}
	return nil
}
