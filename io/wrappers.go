package io

import (
	"errors"
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"io"
)

type luaIOWrapper struct {
	ls  *lua.LState
	obj lua.LValue

	readMethod  *lua.LFunction
	writeMethod *lua.LFunction
	seekMethod  *lua.LFunction
	closeMethod *lua.LFunction
}

//NewLuaIOWrapper creates a new luaIOWrapper atop the lua io object
func NewLuaIOWrapper(L *lua.LState, io lua.LValue) *luaIOWrapper {
	ret := &luaIOWrapper{
		ls:  L,
		obj: io,
	}
	ret.readMethod, _ = L.GetField(io, "read").(*lua.LFunction)
	ret.writeMethod, _ = L.GetField(io, "write").(*lua.LFunction)
	ret.seekMethod, _ = L.GetField(io, "seek").(*lua.LFunction)
	ret.closeMethod, _ = L.GetField(io, "close").(*lua.LFunction)
	return ret
}

//CheckIOWriter tries to cast to UserData and to io.Writer, otherwise it wraps and checks for "write" method
func CheckIOWriter(L *lua.LState, n int) io.Writer {
	any := L.CheckAny(n)
	if ud, ok := any.(*lua.LUserData); ok {
		if writer, ok := ud.Value.(io.Writer); ok {
			return writer
		}
	}
	wrapped := NewLuaIOWrapper(L, any)
	if wrapped.writeMethod == nil {
		L.ArgError(n, "expected writer")
		return nil
	}
	return wrapped
}

//CheckIOReader tries to cast to UserData and to io.Reader, otherwise it wraps and checks for "read" method
func CheckIOReader(L *lua.LState, n int) io.Reader {
	any := L.CheckAny(n)
	if ud, ok := any.(*lua.LUserData); ok {
		if reader, ok := ud.Value.(io.Reader); ok {
			return reader
		}
	}
	wrapped := NewLuaIOWrapper(L, any)
	if wrapped.readMethod == nil {
		L.ArgError(n, "expected reader")
		return nil
	}
	return wrapped
}

func (l *luaIOWrapper) Read(p []byte) (n int, err error) {
	if l.readMethod == nil {
		return 0, errors.New("object does not have read method")
	}
	n = len(p)

	L := l.ls
	L.Push(l.readMethod)
	L.Push(l.obj)
	L.Push(lua.LNumber(n))
	if err = L.PCall(2, 1, nil); err != nil {
		n = 0
		return
	}
	result := L.Get(1)
	L.Pop(L.GetTop())
	if result.Type() == lua.LTNil {
		return 0, io.EOF
	}
	readString := lua.LVAsString(result)
	data := []byte(readString)
	copy(p, data)
	n = len(data)
	return
}

func (l *luaIOWrapper) Write(p []byte) (n int, err error) {
	if l.writeMethod == nil {
		return 0, errors.New("object does not have write method")
	}
	n = len(p)
	L := l.ls
	L.Push(l.writeMethod)
	L.Push(l.obj)
	L.Push(lua.LString(p))
	err = L.PCall(2, 0, nil)
	return
}

func (l *luaIOWrapper) Seek(offset int64, whence int) (int64, error) {
	if l.seekMethod == nil {
		return 0, errors.New("object does not have seek method")
	}
	var luaWhence string
	switch whence {
	case io.SeekStart:
		luaWhence = "set"
	case io.SeekEnd:
		luaWhence = "end"
	case io.SeekCurrent:
		luaWhence = "cur"
	default:
		return 0, fmt.Errorf("unknown whence: %d", whence)
	}

	L := l.ls
	L.Push(l.seekMethod)
	L.Push(l.obj)
	L.Push(lua.LString(luaWhence))
	L.Push(lua.LNumber(offset))
	if err := L.PCall(3, 1, nil); err != nil {
		return 0, err
	}
	ret := L.CheckNumber(1)
	L.Pop(L.GetTop())
	return int64(ret), nil
}

func (l *luaIOWrapper) Close() error {
	if l.closeMethod == nil {
		return errors.New("object does not have close method")
	}
	L := l.ls
	L.Push(l.closeMethod)
	L.Push(l.obj)
	return L.PCall(1, 0, nil)
}
