package zabbix_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"testing"

	http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	zabbix "github.com/vadv/gopher-lua-libs/zabbix"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		zabbix.Preload,
		http.Preload,
		inspect.Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
