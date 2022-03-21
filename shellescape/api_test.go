package shellescape

import (
	"github.com/vadv/gopher-lua-libs/inspect"
	"github.com/vadv/gopher-lua-libs/tests"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestApi(t *testing.T) {
	preloadForTest := func(L *lua.LState) {
		Preload(L)
		inspect.Preload(L)
	}
	tests.RunLuaTestFile(t, preloadForTest, "test/test_api.lua")
}
