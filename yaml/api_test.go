package yaml

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestApi(t *testing.T) {
	// Setup the VM
	L := lua.NewState()
	defer L.Close()
	Preload(L)

	// Load the test cases from file
	err := L.DoFile("./test/test_api.lua")
	require.NoError(t, err)
	test := L.CheckTable(1)
	L.Pop(1)

	// For each method on the returned test object, invoke it safely with PCall.
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

	// Ensure we ran non-zero tests.
	assert.NotEqual(t, 0, testCount, "test should not be empty")
}
