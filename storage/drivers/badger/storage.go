// this storage based on github.com/dgraph-io/badger
package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

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
		result.Lock()
		defer result.Unlock()
		result.usageCounter++
		return result, nil
	}

	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path
	opts.TableLoadingMode = badger_options.FileIO
	opts.ValueLogLoadingMode = badger_options.FileIO
	opts.SyncWrites = false
	opts.NumCompactors = 1
	opts.MaxTableSize = 1024 * 1024
	if sizeStr := os.Getenv(`BADGER_MAX_TABLE_SIZE_MB`); sizeStr != `` {
		size, err := strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bad value for BADGER_MAX_TABLE_SIZE_MB: %s", err.Error())
		}
		opts.MaxTableSize = size * 1024 * 1024
	}
	opts.Truncate = true

	badgerDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	s := &Storage{DB: badgerDB, path: path}
	go s.gc()
	log.Printf("[INFO] new badger storage [%p-%s]\n", s, s.path)
	s.usageCounter++
	listOfStorages.list[path] = s
	return s, nil
}

func (s *Storage) gc() {
	for {
		time.Sleep(5 * time.Minute)
		now := time.Now()
		s.Lock()
		if s.usageCounter == 0 {
			return
		}
		err := s.DB.RunValueLogGC(0.7)
		s.Unlock()
		if err != nil {
			log.Printf("[ERROR] [%p-%s] while running gc: %v\n", s, s.path, err.Error())
		} else {
			log.Printf("[INFO] [%p-%s] gc completed, execution time: %v\n", s, s.path, time.Now().Sub(now).Seconds())
		}
	}
}

func (s *Storage) Sync() error {
	return nil
}

func (s *Storage) Close() error {
	listOfStorages.Lock()
	defer listOfStorages.Unlock()
	s.Lock()
	defer s.Unlock()
	s.usageCounter--
	if s.usageCounter == 0 {
		log.Printf("[INFO] close unused badger storage [%p-%s]\n", s, s.path)
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
			result = append(result, string(k))
		}
		return nil
	})
	return result, err
}

func (s *Storage) Dump(L *lua.LState) (map[string]lua.LValue, error) {
	return nil, fmt.Errorf("unsupported")
}
