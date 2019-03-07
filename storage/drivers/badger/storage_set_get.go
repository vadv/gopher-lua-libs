package storage

import (
	"time"

	lua_json "github.com/vadv/gopher-lua-libs/json"

	lua "github.com/yuin/gopher-lua"
	badger "gopkg.in/dgraph-io/badger.v1"
)

func (s *Storage) Set(key string, value lua.LValue, ttl int64) error {
	data, err := lua_json.ValueEncode(value)
	if err != nil {
		return err
	}
	err = s.DB.Update(func(txn *badger.Txn) error {
		if ttl > 0 {
			return txn.SetWithTTL([]byte(key), data, time.Duration(ttl)*time.Second)
		} else {
			return txn.Set([]byte(key), data)
		}
	})
	return err
}

func (s *Storage) Get(key string, L *lua.LState) (lua.LValue, bool, error) {
	var result lua.LValue
	err := s.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		data, err := item.Value()
		if err != nil {
			return err
		}
		value, err := lua_json.ValueDecode(L, data)
		if err != nil {
			return err
		}
		result = value
		return nil
	})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return lua.LNil, false, nil
		}
		return lua.LNil, false, err
	}
	return result, true, nil
}
