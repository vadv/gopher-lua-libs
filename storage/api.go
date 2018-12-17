// Package storage implements persist storage with ttl for to save and share data between differents lua.LState.
package storage

import (
	lua "github.com/yuin/gopher-lua"
)

// New(): lua storage.new(filename) returns (storage_ud, err)
func New(L *lua.LState) int {
	filename := L.CheckString(1)
	s, err := newStorage(filename)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = s
	L.SetMetatable(ud, L.GetTypeMetatable("storage_ud"))
	L.Push(ud)
	return 1
}

// Set(): lua storage_ud:set(key, value, ttl) return err
func Set(L *lua.LState) int {
	s := checkStorage(L, 1)
	key := L.CheckString(2)
	value := L.CheckAny(3)
	ttl := int64(60)
	if L.GetTop() > 3 {
		luaTTL := L.CheckAny(4)
		switch luaTTL.(type) {
		case *lua.LNilType:
			ttl = 1000000000000 // max tll in second :)
		case lua.LNumber:
			ttl = L.CheckInt64(4)
		default:
			L.ArgError(4, "must be integer or nil")
		}
	}
	err := s.set(key, value, ttl)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Get(): lua storage_ud:set(key) returns (value, bool, err)
func Get(L *lua.LState) int {
	s := checkStorage(L, 1)
	key := L.CheckString(2)
	value, found, err := s.get(key, L)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 3
	}
	L.Push(value)
	L.Push(lua.LBool(found))
	return 2
}

// Sync(): lua storage_ud:sync() return err
func Sync(L *lua.LState) int {
	s := checkStorage(L, 1)
	err := s.sync()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Close(): lua storage_ud:close() return err
func Close(L *lua.LState) int {
	s := checkStorage(L, 1)
	err := s.close()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
