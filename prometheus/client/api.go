package prometheus_client

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	lua "github.com/yuin/gopher-lua"
)

type luaPrometheusClient struct {
	addr string
	stop chan bool
}

func checkPrometheusClient(L *lua.LState, n int) *luaPrometheusClient {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaPrometheusClient); ok {
		return v
	}
	L.ArgError(n, "prometheus_client_ud expected")
	return nil
}

// Register(string): return (prometheus_client_ud, err)
func Register(L *lua.LState) int {
	addr := L.CheckString(1)
	ud := L.NewUserData()
	ud.Value = &luaPrometheusClient{addr: addr, stop: make(chan bool, 1)}
	L.SetMetatable(ud, L.GetTypeMetatable(`prometheus_client_ud`))
	L.Push(ud)
	return 1
}

// Start prometheus_client_ud
func Start(L *lua.LState) int {
	pp := checkPrometheusClient(L, 1)
	go func() {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		h := &http.Server{Addr: pp.addr, Handler: m}

		// register shutdown
		go func() {
			ctx := L.Context()
			if ctx != nil {
				select {
				case <-ctx.Done():
					pp.stop <- true
				}
			}
		}()

		// start listen
		go func() {
			if err := h.ListenAndServe(); err != nil {
				return
			}
		}()

		// wait shutdown
		<-pp.stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		h.Shutdown(ctx)
	}()
	return 0
}

// Stop prometheus_client_ud stop
func Stop(L *lua.LState) int {
	pp := checkPrometheusClient(L, 1)
	pp.stop <- true
	return 0
}
