// Package tac implements file scanner (from end to up) functionality for lua.
package tac

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

type luaTac struct {
	filename string
	fd       *os.File
	scanner  *tacScanner
}

func checkTac(L *lua.LState, n int) *luaTac {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*luaTac); ok {
		return v
	}
	L.ArgError(1, "tac expected")
	return nil
}

func (t *luaTac) open() error {
	fd, err := os.Open(t.filename)
	if err != nil {
		return err
	}
	t.fd = fd
	t.scanner = newTacScanner(fd)
	return nil
}

// Open lua tac.open(filename) open filename for tac scan returns (tac_ud, err)
func Open(L *lua.LState) int {
	t := &luaTac{filename: L.CheckString(1)}
	if err := t.open(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = t
	L.SetMetatable(ud, L.GetTypeMetatable("tac_ud"))
	L.Push(ud)
	return 1
}

// Line lua tac_ud:line() return next upper line: string or nil
func Line(L *lua.LState) int {
	t := checkTac(L, 1)
	if t.scanner == nil {
		L.RaiseError("tac not initialized")
		return 0
	}
	if t.scanner.scan() {
		text := t.scanner.text()
		L.Push(lua.LString(text))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

// Close lua tac_ud:close() close current file for tac
func Close(L *lua.LState) int {
	t := checkTac(L, 1)
	if t.fd == nil {
		L.RaiseError("tac not initialized")
		return 0
	}
	t.fd.Close()
	t = nil
	return 0
}
