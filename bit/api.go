// Package strings implements golang package montanaflynn/stats functionality for lua.
package stats

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

// ShiftLeft
func Bitwise(kind op) lua.LGFunction {
	return func(l *lua.LState) int {
		if kind > rs {
			l.Push(lua.LString("invalid type of operation"))
			return 1
		}
		val1, val2, err := prepareParams(l)
		if err != nil {
			l.Push(lua.LString(err.Error()))
			return 1
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

func Not(l *lua.LState) int {
	val, err := intToU32(l.CheckInt(1))
	if err != nil {
		l.Push(lua.LString(err.Error()))
		return 1
	}
	l.Push(lua.LNumber(^val))
	return 1
}

func prepareParams(l *lua.LState) (val, pos uint32, err error) {
	val, err = intToU32(l.CheckInt(1))
	if err != nil {
		return 0, 0, err
	}
	pos, err = intToU32(l.CheckInt(2))
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
