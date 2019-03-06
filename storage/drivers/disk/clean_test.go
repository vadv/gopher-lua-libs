package storage

import (
	"os"
	"testing"
	"time"
)

func TestClean(t *testing.T) {
	dir := "/xxxx/yyyy/zzzz/"
	if err := os.MkdirAll("./test"+dir, 0750); err != nil {
		t.Fatal(err)
	}
	s := &Storage{path: "./test/"}
	h := newHeaderInfo("xxx", 1)
	headerFile := "./test" + dir + "file" + headerExt
	h.write(headerFile)
	time.Sleep(time.Second * 2)
	s.cleanRoutine()
	if _, err := os.Stat(headerFile); !os.IsNotExist(err) {
		t.Fatal("must be deleted")
	}
	s.cleanRoutine()
	s.cleanRoutine()
	s.cleanRoutine()
}
