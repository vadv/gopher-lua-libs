package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	headerExt = `.header`
)

type headerInfo struct {
	key string
	ttl int64
}

func parseHeader(path string) (*headerInfo, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(data), "\n")
	if len(list) != 2 {
		return nil, fmt.Errorf("unsupported version")
	}
	result := &headerInfo{key: list[0]}
	ttl, err := strconv.ParseInt(list[1], 10, 64)
	if err != nil {
		return nil, err
	}
	result.ttl = ttl
	return result, nil
}

func newHeaderInfo(key string, ttl int64) *headerInfo {
	h := &headerInfo{key: key}
	h.ttl = time.Now().UnixNano() + (ttl * 1000000000)
	return h
}

func (h *headerInfo) hasValidTTL() bool {
	return h.ttl > time.Now().UnixNano()
}

func (h *headerInfo) write(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return err
	}
	data := fmt.Sprintf("%s\n%d", h.key, h.ttl)
	return ioutil.WriteFile(path, []byte(data), 0640)
}
