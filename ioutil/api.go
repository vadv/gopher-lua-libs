// Package ioutil implements golang package ioutil functionality for lua.
package ioutil

import (
	"io/ioutil"

	lua "github.com/yuin/gopher-lua"
)

// ReadFile lua ioutil.read_file(filepath) reads the file named by filename and returns the contents, returns (string,error)
func ReadFile(L *lua.LState) int {
	filename := L.CheckString(1)
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		L.Push(lua.LString(data))
		return 1
	} else {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
}

// WriteFile lua ioutil.write_file(filepath, data) reads the file named by filename and returns the contents, returns (string,error)
func WriteFile(L *lua.LState) int {
	filename := L.CheckString(1)
	data := L.CheckString(2)
	err := ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
