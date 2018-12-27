// Package goos implements golang package os functionality for lua.
package goos

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

// Stat(): lua os.stat(filename) returns (table, err)
func Stat(L *lua.LState) int {
	filename := L.CheckString(1)
	stat, err := os.Stat(filename)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	result.RawSetString(`is_dir`, lua.LBool(stat.IsDir()))
	result.RawSetString(`size`, lua.LNumber(stat.Size()))
	result.RawSetString(`mod_time`, lua.LNumber(stat.ModTime().Unix()))
	result.RawSetString(`mode`, lua.LString(stat.Mode().String()))
	L.Push(result)
	return 1
}
