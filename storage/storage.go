package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	lua_json "github.com/vadv/gopher-lua-libs/json"

	lua "github.com/yuin/gopher-lua"
)

var listOfStorages = &listStorages{list: make(map[string]*storage)}

type listStorages struct {
	sync.Mutex
	list map[string]*storage
}

type storage struct {
	sync.Mutex
	filename string
	Data     map[string]*storageValue `json:"data"`
	running  bool
}

func checkStorage(L *lua.LState, n int) *storage {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*storage); ok {
		return v
	}
	L.ArgError(n, "storage_ud excepted")
	return nil
}

type storageValue struct {
	Value      []byte `json:"value"`        // json value
	MaxValidAt int64  `json:"max_valid_at"` // unix nano
}

func (sv *storageValue) valid() bool {
	return sv.MaxValidAt > time.Now().UnixNano()
}

func newStorage(filename string) (*storage, error) {

	listOfStorages.Lock()
	defer listOfStorages.Unlock()

	if result, ok := listOfStorages.list[filename]; ok {
		return result, nil
	}

	s := &storage{Data: make(map[string]*storageValue, 0)}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// create
		dst, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}
		dst.Close()
	} else {
		// read && decode
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, s); err != nil {
			return nil, err
		}
	}
	s.filename = filename
	s.running = true
	listOfStorages.list[filename] = s
	go s.loop()
	return s, s.sync()
}

func (s *storage) sync() error {
	s.Lock()
	defer s.Unlock()
	tmpFilename := s.filename + ".tmp"
	// clean
	newData := make(map[string]*storageValue, 0)
	for k, v := range s.Data {
		if v.valid() {
			newData[k] = v
		}
	}
	s.Data = newData
	// clean end
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(tmpFilename, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmpFilename, s.filename)
}

func (s *storage) close() error {
	listOfStorages.Lock()
	defer listOfStorages.Unlock()
	if err := s.sync(); err != nil {
		return err
	}
	delete(listOfStorages.list, s.filename)
	s.running = false
	return nil
}

func (s *storage) loop() {
	for {
		time.Sleep(60 * time.Second)
		if !s.running {
			return
		}
		if err := s.sync(); err != nil {
			log.Printf("[ERROR] scheduler [%p] sync save %s: %s\n", s, s.filename, err.Error())
		}
	}
}

func (s *storage) keys() []string {
	result := []string{}
	s.Lock()
	defer s.Unlock()
	for k, _ := range s.Data {
		result = append(result, k)
	}
	return result
}

func (s *storage) dump(L *lua.LState) (map[string]lua.LValue, error) {
	result := make(map[string]lua.LValue, 0)
	s.Lock()
	defer s.Unlock()
	for k, v := range s.Data {
		if v.valid() {
			value, err := lua_json.ValueDecode(L, v.Value)
			if err != nil {
				return nil, err
			}
			result[k] = value
		}
	}
	return result, nil
}
