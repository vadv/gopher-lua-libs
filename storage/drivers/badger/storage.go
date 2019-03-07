// this storage based on github.com/dgraph-io/badger
package storage

import (
	"fmt"
	"sync"

	interfaces "github.com/vadv/gopher-lua-libs/storage/drivers/interfaces"

	badger "github.com/dgraph-io/badger"
	badger_options "github.com/dgraph-io/badger/options"

	lua "github.com/yuin/gopher-lua"
)

var listOfStorages = &listStorages{list: make(map[string]*Storage)}

type listStorages struct {
	sync.Mutex
	list map[string]*Storage
}

type Storage struct {
	*badger.DB
	path         string
	usageCounter int
}

func (st *Storage) New(path string) (interfaces.Driver, error) {

	listOfStorages.Lock()
	defer listOfStorages.Unlock()

	if result, ok := listOfStorages.list[path]; ok {
		result.usageCounter++
		return result, nil
	}

	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path
	opts.TableLoadingMode = badger_options.MemoryMap
	opts.Truncate = true

	badgerDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	s := &Storage{DB: badgerDB}
	s.usageCounter++
	listOfStorages.list[path] = s
	return s, nil
}

func (s *Storage) Sync() error {
	return nil
}

func (s *Storage) Close() error {
	listOfStorages.Lock()
	defer listOfStorages.Unlock()
	s.usageCounter--
	if s.usageCounter == 0 {
		return s.DB.Close()
	}
	return nil
}

func (s *Storage) Keys() ([]string, error) {
	result := []string{}
	err := s.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			result = append(result, fmt.Sprintf("%s", k))
		}
		return nil
	})
	return result, err
}

func (s *Storage) Dump(L *lua.LState) (map[string]lua.LValue, error) {
	return nil, fmt.Errorf("unsupported")
}
