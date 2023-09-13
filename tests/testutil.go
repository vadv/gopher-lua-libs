package tests

import (
	_ "embed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

//go:embed suite.lua
var lua_suite string

//go:embed assertions.lua
var lua_assertions string

//go:embed assert.lua
var lua_assert string

//go:embed require.lua
var lua_require string

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
	args := []interface{}{"\n", where}
	top := L.GetTop()
	for i := 2; i <= top; i++ {
		args = append(args, L.Get(i))
	}
	t.Log(args...)
	return 0
}

func tLogHelper(L *lua.LState) int {
	t := checkT(L, 1)
	level := L.CheckInt(2)
	where := L.Where(level)
	args := []interface{}{"\n", where}
	top := L.GetTop()
	for i := 3; i <= top; i++ {
		args = append(args, L.Get(i))
	}
	t.Log(args...)
	return 0
}

func tLogf(L *lua.LState) int {
	t := checkT(L, 1)
	format := "\n%s " + L.CheckString(2)
	where := L.Where(1)
	args := []interface{}{where}
	top := L.GetTop()
	for i := 3; i <= top; i++ {
		args = append(args, L.Get(i))
	}
	t.Logf(format, args...)
	return 0
}

func tLogHelperf(L *lua.LState) int {
	t := checkT(L, 1)
	level := L.CheckInt(2)
	format := "\n%s " + L.CheckString(3)
	where := L.Where(level)
	args := []interface{}{where}
	top := L.GetTop()
	for i := 4; i <= top; i++ {
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

func tFail(L *lua.LState) int {
	t := checkT(L, 1)
	t.Fail()
	return 0
}

func tFailNow(L *lua.LState) int {
	t := checkT(L, 1)
	t.FailNow()
	return 0
}

func tFailed(L *lua.LState) int {
	t := checkT(L, 1)
	L.Push(lua.LBool(t.Failed()))
	return 1
}

func tError(L *lua.LState) int {
	tLog(L)
	t := checkT(L, 1)
	t.Fail()
	return 0
}

func tErrorf(L *lua.LState) int {
	tLogf(L)
	t := checkT(L, 1)
	t.Fail()
	return 0
}

func tFatal(L *lua.LState) int {
	tLog(L)
	t := checkT(L, 1)
	t.FailNow()
	return 0
}

func tFatalf(L *lua.LState) int {
	tLogf(L)
	t := checkT(L, 1)
	t.FailNow()
	return 0
}

func tTempDir(L *lua.LState) int {
	t := checkT(L, 1)
	// TODO(scr): When the minimal version supported has this on the *testing.T object, remove this shim
	//L.Push(lua.LString(t.TempDir()))
	tempDir, err := ioutil.TempDir(os.TempDir(), "test.tempDir*")
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = os.RemoveAll(tempDir)
	})

	L.Push(lua.LString(tempDir))
	return 1
}

func registerTType(L *lua.LState) {
	mt := L.NewTypeMetatable(TType)
	index := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"Error":      tError,
		"Errorf":     tErrorf,
		"Fail":       tFail,
		"FailNow":    tFailNow,
		"Failed":     tFailed,
		"Fatal":      tFatal,
		"Fatalf":     tFatalf,
		"Log":        tLog,
		"LogHelper":  tLogHelper,
		"LogHelperf": tLogHelperf,
		"Logf":       tLogf,
		"Run":        tRun,
		"Skip":       tSkip,
		"Skipf":      tSkipf,
		"TempDir":    tTempDir,
	})
	L.SetField(mt, "__index", index)
	L.SetGlobal(TType, mt)
}

func LoadSuite(L *lua.LState) int {
	if err := L.DoString(lua_suite); err != nil {
		L.RaiseError(err.Error())
	}
	return 1
}

func PreloadSuite(L *lua.LState) {
	L.PreloadModule("suite", LoadSuite)
}

func LoadAssertions(L *lua.LState) int {
	if err := L.DoString(lua_assertions); err != nil {
		L.RaiseError(err.Error())
	}
	return 1
}

func PreloadAssertions(L *lua.LState) {
	L.PreloadModule("assertions", LoadAssertions)
}

func LoadAssert(L *lua.LState) int {
	if err := L.DoString(lua_assert); err != nil {
		L.RaiseError(err.Error())
	}
	return 1
}

func PreloadAssert(L *lua.LState) {
	L.PreloadModule("assert", LoadAssert)
}

func LoadRequre(L *lua.LState) int {
	if err := L.DoString(lua_require); err != nil {
		L.RaiseError(err.Error())
	}
	return 1
}

func PreloadRequire(L *lua.LState) {
	L.PreloadModule("require", LoadRequre)
}

// RunLuaTestFile fires up a new state, registers the *testing.T and invokes all methods starting with Test.
// This allows the lua test files to operate similar to go tests - see shellescape/test/test_api.lua
func RunLuaTestFile(t *testing.T, preload PreloadFunc, filename string) (numTests int) {
	L := lua.NewState()
	t.Cleanup(L.Close)

	registerTType(L)
	PreloadSuite(L)
	PreloadAssertions(L)
	PreloadAssert(L)
	PreloadRequire(L)
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

// SeveralPreloadFuncs combines several PreloadFuncs to one such as when tests want to preload theirs + inspect
func SeveralPreloadFuncs(preloadFuncs ...PreloadFunc) PreloadFunc {
	return func(L *lua.LState) {
		for _, preloadFunc := range preloadFuncs {
			preloadFunc(L)
		}
	}
}
