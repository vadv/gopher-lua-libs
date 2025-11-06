package goos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"

	runtime "github.com/vadv/gopher-lua-libs/runtime"
)

func TestApi(t *testing.T) {
	os.Setenv("ENV_VAR", "TEST=1")
	defer os.Unsetenv("ENV_VAR")

	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("EMPTY_VAR")

	preload := tests.SeveralPreloadFuncs(
		runtime.Preload,
		Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
