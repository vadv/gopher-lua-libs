// Package time implements golang package time functionality for lua.
package time

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

// Unix(): lua time.unix() returns unix timestamp in seconds (float)
func Unix(L *lua.LState) int {
	now := float64(time.Now().UnixNano()) / float64(time.Second)
	L.Push(lua.LNumber(now))
	return 1
}

// UnixNano(): lua time.unix_nano() returns unix timestamp in nano seconds
func UnixNano(L *lua.LState) int {
	L.Push(lua.LNumber(time.Now().UnixNano()))
	return 1
}

// Sleep(): lua time.sleep(number) port of go time.Sleep(int64)
func Sleep(L *lua.LState) int {
	val := L.CheckNumber(1)
	time.Sleep(time.Duration(val) * time.Second)
	return 0
}

// Parse(): lua time.parse(value, layout) returns (number, error)
func Parse(L *lua.LState) int {
	layout, value := L.CheckString(2), L.CheckString(1)
	result, err := time.Parse(layout, value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	resultFloat := float64(result.UTC().UnixNano()) / float64(time.Second)
	L.Push(lua.LNumber(resultFloat))
	return 1
}
