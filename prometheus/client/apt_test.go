package prometheus_client_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/http"
	prometheus "github.com/vadv/gopher-lua-libs/prometheus/client"
	"github.com/vadv/gopher-lua-libs/strings"
	"github.com/vadv/gopher-lua-libs/tests"
	"github.com/vadv/gopher-lua-libs/time"
	"testing"
)

func TestApi(t *testing.T) {
	preload := tests.SeveralPreloadFuncs(
		prometheus.Preload,
		http.Preload,
		strings.Preload,
		time.Preload,
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
