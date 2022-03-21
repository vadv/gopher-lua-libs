package shellescape

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadv/gopher-lua-libs/inspect"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	L := lua.NewState()
	Preload(L)
	inspect.Preload(L)

	// Load the test cases from file
	require.NoError(t, L.DoFile("test/test_api.lua"))
	require.Equal(t, 1, L.GetTop())
	test := L.CheckTable(1)
	L.Pop(1)

	testCount := 0
	test.ForEach(func(key lua.LValue, value lua.LValue) {
		if value.Type() != lua.LTFunction {
			return
		}
		testCount++
		t.Run(lua.LVAsString(key), func(t *testing.T) {
			L.Push(value)
			L.Push(test)
			require.NoError(t, L.PCall(1, 0, nil))
		})
	})

	assert.Greater(t, testCount, 0, "No tests were run")
}
