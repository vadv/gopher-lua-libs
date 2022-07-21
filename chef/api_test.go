package chef_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"

	chef "github.com/vadv/gopher-lua-libs/chef"
	http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		chef.Preload,
		http.Preload,
		inspect.Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
