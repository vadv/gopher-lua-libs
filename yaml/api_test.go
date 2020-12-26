package yaml

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	Preload(L)
	err := L.DoFile("./test/test_api.lua")
	require.NoError(t, err)
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
	assert.NotEqual(t, 0, testCount, "test should not be empty")
}
