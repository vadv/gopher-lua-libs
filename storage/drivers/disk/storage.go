// this storage for projects that store a lot of data and save memory
package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	interfaces "github.com/vadv/gopher-lua-libs/storage/drivers/interfaces"

	lua "github.com/yuin/gopher-lua"
)

var listOfStorages = &listStorages{list: make(map[string]*Storage)}

type listStorages struct {
	sync.Mutex
	list map[string]*Storage
}

type Storage struct {
	sync.Mutex
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

	// check
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(path, 0750)
		} else {
			return nil, err
		}
	} else {
		if !stat.IsDir() {
			return nil, fmt.Errorf("must be directory")
		}
	}

	s := &Storage{path: path}
	s.usageCounter++
	listOfStorages.list[path] = s
	go s.loop()
	return s, nil
}

func (s *Storage) Sync() error {
	return nil
}

func (s *Storage) Close() error {
	listOfStorages.Lock()
	defer listOfStorages.Unlock()
	s.usageCounter--
	return nil
}

// cleaner and closer
func (s *Storage) loop() {
	for {
		time.Sleep(5 * time.Minute)
		s.cleanRoutine()
	}
}

func (s *Storage) Keys() ([]string, error) {
	headerGlobPatern := filepath.Join(s.path,
		fmt.Sprintf(filePath, "*", "*", "*")+headerExt)
	files, err := filepath.Glob(headerGlobPatern)
	if err != nil {
		return nil, err
	}
	results := []string{}
	for _, file := range files {
		h, err := parseHeader(file)
		if err != nil {
			return nil, err
		}
		if h.hasValidTTL() {
			results = append(results, h.key)
		}
	}
	return results, nil
}

func (s *Storage) Dump(L *lua.LState) (map[string]lua.LValue, error) {
	return nil, fmt.Errorf("unsupported")
}
