package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"strings"
	"testing"
)

type PreloadFunc func(L *lua.LState)

const (
	TType = "testing.T"
)

func tLua(L *lua.LState, t *testing.T) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = t
	L.SetMetatable(ud, L.GetTypeMetatable(TType))
	return ud
}

func checkT(L *lua.LState, n int) *testing.T {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*testing.T); ok {
		return v
	}
	L.ArgError(n, "testing.T expected")
	return nil
}

func tRun(L *lua.LState) int {
	t := checkT(L, 1)
	name := L.CheckString(2)
	function := L.CheckFunction(3)
	L.Pop(L.GetTop())

	t.Run(name, func(t *testing.T) {
		L.Push(function)
		L.Push(tLua(L, t))
		assert.NoError(t, L.PCall(1, 0, nil))
	})

	return 0
}

func tLog(L *lua.LState) int {
	t := checkT(L, 1)
	where := L.Where(1)
	args := []interface{}{where}
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		args = append(args, L.Get(i))
	}
	t.Log(args...)
	return 0
}

func tLogf(L *lua.LState) int {
	t := checkT(L, 1)
	format := "%s " + L.CheckString(2)
	where := L.Where(1)
	args := []interface{}{where}
	top := L.GetTop()
	for i := 3; i <= top; i++ {
		args = append(args, L.Get(i))
	}
	t.Logf(format, args...)
	return 0
}

func tSkip(L *lua.LState) int {
	t := checkT(L, 1)
	var args []interface{}
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		args = append(args, L.Get(i).String())
	}
	t.Skip(args...)
	return 0
}

func tSkipf(L *lua.LState) int {
	t := checkT(L, 1)
	format := L.CheckString(2)
	var args []interface{}
	top := L.GetTop()
	for i := 3; i <= top; i++ {
		args = append(args, L.Get(i).String())
	}
	t.Skipf(format, args...)
	return 0
}

func registerTType(L *lua.LState) {
	mt := L.NewTypeMetatable(TType)
	index := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"Run":   tRun,
		"Log":   tLog,
		"Logf":  tLogf,
		"Skip":  tSkip,
		"Skipf": tSkipf,
	})
	L.SetField(mt, "__index", index)
	L.SetGlobal(TType, mt)
}

//RunLuaTestFile fires up a new state, registers the *testing.T and invokes all methods starting with Test.
// This allows the lua test files to operate similar to go tests - see shellescape/test/test_api.lua
func RunLuaTestFile(t *testing.T, preload PreloadFunc, filename string) (numTests int) {
	L := lua.NewState()
	t.Cleanup(L.Close)

	registerTType(L)
	require.NotNil(t, preload)
	preload(L)
	L.SetGlobal("t", tLua(L, t))

	require.NoError(t, L.DoFile(filename))
	L.G.Global.ForEach(func(key lua.LValue, value lua.LValue) {
		keyStr := lua.LVAsString(key)
		if strings.HasPrefix(keyStr, "Test") && value.Type() == lua.LTFunction {
			t.Run(keyStr, func(t *testing.T) {
				numTests++
				L.Push(value)
				L.Push(tLua(L, t))
				assert.NoError(t, L.PCall(1, 0, nil))
			})
		}
	})
	return
}

//SeveralPreloadFuncs combines several PreloadFuncs to one such as when tests want to preload theirs + inspect
func SeveralPreloadFuncs(preloadFuncs ...PreloadFunc) PreloadFunc {
	return func(L *lua.LState) {
		for _, preloadFunc := range preloadFuncs {
			preloadFunc(L)
		}
	}
}
