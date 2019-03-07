package storage

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	lua_json "github.com/vadv/gopher-lua-libs/json"

	lua "github.com/yuin/gopher-lua"
)

var (
	separator = string(filepath.Separator)
	filePath  = `%s` + separator + `%s` + separator + `%s`
)

func (s *Storage) getFilePath(key string) string {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	hash1, hash2, hash3 := hash[0:4], hash[5:8], hash[9:]
	return filepath.Join(s.path, fmt.Sprintf(filePath, hash1, hash2, hash3))
}

func (s *Storage) get(key string) ([]byte, bool, error) {
	filePath := s.getFilePath(key)
	if _, err := os.Stat(filePath + headerExt); err != nil && os.IsNotExist(err) {
		return nil, false, nil
	}
	header, err := parseHeader(filePath + headerExt)
	if err != nil {
		return nil, true, err
	}
	if key != header.key {
		return nil, true, fmt.Errorf("bad key in header")
	}
	if !header.hasValidTTL() {
		if err := os.RemoveAll(filePath + headerExt); err != nil {
			return nil, false, err
		}
		os.RemoveAll(filePath)
		os.Remove(filepath.Dir(filePath))
		os.Remove(filepath.Dir(filepath.Dir(filePath)))
		return nil, false, nil
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, true, err
	}
	return data, true, nil
}

func (s *Storage) Set(key string, value lua.LValue, ttl int64) error {
	s.Lock()
	defer s.Unlock()
	filePath := s.getFilePath(key)
	data, err := lua_json.ValueEncode(value)
	if err != nil {
		return err
	}
	header := newHeaderInfo(key, ttl)
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, data, 0640); err != nil {
		return err
	}
	if err := header.write(filePath + headerExt); err != nil {
		return err
	}
	time.Sleep(10 * time.Millisecond)
	return nil
}

func (s *Storage) Get(key string, L *lua.LState) (lua.LValue, bool, error) {
	s.Lock()
	defer s.Unlock()
	data, found, err := s.get(key)
	if !found {
		return lua.LNil, false, nil
	}
	if err != nil {
		return lua.LNil, true, err
	}
	value, err := lua_json.ValueDecode(L, data)
	if err != nil {
		return lua.LNil, false, err
	}
	return value, true, nil
}
