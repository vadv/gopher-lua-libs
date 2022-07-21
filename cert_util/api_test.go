package cert_util

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/vadv/gopher-lua-libs/tests"
	"io"
	"net/http"
	"testing"
	"time"
)

func runHttps(addr string, handler http.Handler) *http.Server {
	server := &http.Server{Addr: addr, Handler: handler}
	go func() {
		_ = server.ListenAndServeTLS("./test/cert.pem", "./test/key.pem")
	}()
	return server
}

func httpRouterGet(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "OK")
}

func TestApi(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", httpRouterGet)
	server := runHttps(":1443", mux)
	t.Cleanup(func() {
		_ = server.Shutdown(context.Background())
	})
	time.Sleep(time.Second)

	assert.NotZero(t, tests.RunLuaTestFile(t, Preload, "./test/test_api.lua"))
}
