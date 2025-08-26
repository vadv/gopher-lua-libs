// Package bit implements Go bitwise operations functionality for Lua.
package bit

import (
	"fmt"
	"math"

	lua "github.com/yuin/gopher-lua"
)

type op uint

const (
	and op = iota
	or
	not
	xor
	ls
	rs
)

// Bitwise returns a Lua function used for bitwise operations.
func Bitwise(kind op) lua.LGFunction {
	return func(l *lua.LState) int {
		if kind > rs {
			l.RaiseError("unsupported operation type")
			return 0
		}
		val1, val2, err := prepareParams(l)
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
			return 2
		}
		var ret uint32
		switch kind {
		case and:
			ret = val1 & val2
		case or:
			ret = val1 | val2
		case xor:
			ret = val1 ^ val2
		case ls:
			ret = val1 << val2
		case rs:
			ret = val1 >> val2
		}
		l.Push(lua.LNumber(ret))
		return 1
	}
}

// Not implements bitwise not.
func Not(l *lua.LState) int {
	val, err := intToU32(l.CheckInt(1))
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}
	l.Push(lua.LNumber(^val))
	return 1
}

func prepareParams(l *lua.LState) (val1, val2 uint32, err error) {
	val1, err = intToU32(l.CheckInt(1))
	if err != nil {
		return 0, 0, err
	}
	val2, err = intToU32(l.CheckInt(2))
	if err != nil {
		return 0, 0, err
	}
	return
}

func intToU32(i int) (uint32, error) {
	if i < 0 {
		return 0, fmt.Errorf("cannot convert negative int %d to uint32", i)
	}
	if i > math.MaxUint32 {
		return 0, fmt.Errorf("int %d overflows uint32", i)
	}
	return uint32(i), nil
}
