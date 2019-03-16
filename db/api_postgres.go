package db

import (
	"database/sql"
	"sync"

	_ "github.com/lib/pq"
)

type luaPG struct {
	config *dbConfig
	sync.Mutex
	db *sql.DB
}

func init() {
	RegisterDriver(`postgres`, &luaPG{})
}

var (
	sharedPG     = make(map[string]*luaPG, 0)
	sharedPGLock = &sync.Mutex{}
)

func (pg *luaPG) constructor(config *dbConfig) (luaDB, error) {

	sharedPGLock.Lock()
	defer sharedPGLock.Unlock()

	if config.sharedMode {
		result, ok := sharedPG[config.connString]
		if ok {
			return result, nil
		}
	}

	db, err := sql.Open(`postgres`, config.connString)
	if err != nil {
		return nil, err
	}
	result := &luaPG{config: config}
	db.SetMaxIdleConns(config.maxOpenConns)
	db.SetMaxOpenConns(config.maxOpenConns)
	result.db = db

	if config.sharedMode {
		sharedPG[config.connString] = result
	}

	return result, nil
}

func (pg *luaPG) getTXOptions() *sql.TxOptions {
	return &sql.TxOptions{ReadOnly: pg.config.readOnly}
}

func (pg *luaPG) getDB() *sql.DB {
	pg.Lock()
	defer pg.Unlock()
	return pg.db
}

func (pg *luaPG) closeDB() error {
	pg.Lock()
	defer pg.Unlock()
	err := pg.db.Close()
	if err != nil {
		return err
	}
	if pg.config.sharedMode {
		sharedPGLock.Lock()
		delete(sharedPG, pg.config.connString)
		sharedPGLock.Unlock()
	}
	return nil
}
