package telegram_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"

	http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	telegram "github.com/vadv/gopher-lua-libs/telegram"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		telegram.Preload,
		http.Preload,
		inspect.Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
