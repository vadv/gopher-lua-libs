// Package pprof implements golang package pprof functionality.
package pprof

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type luaPprof struct {
	addr string
	stop chan bool
}

func checkPprof(L *lua.LState, n int) *luaPprof {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaPprof); ok {
		return v
	}
	L.ArgError(n, "pprof_ud expected")
	return nil
}

// Register(string): return (pprof_ud, err)
func Register(L *lua.LState) int {
	addr := L.CheckString(1)
	ud := L.NewUserData()
	ud.Value = &luaPprof{addr: addr, stop: make(chan bool, 1)}
	L.SetMetatable(ud, L.GetTypeMetatable(`pprof_ud`))
	L.Push(ud)
	return 1
}

// Enable start pprof
func Enable(L *lua.LState) int {
	pp := checkPprof(L, 1)
	go func() {
		h := &http.Server{Addr: pp.addr}
		go func() {
			if err := h.ListenAndServe(); err != nil {
				return
			}
		}()
		<-pp.stop
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		h.Shutdown(ctx)
	}()
	return 0
}

// Disable pprof stop
func Disable(L *lua.LState) int {
	pp := checkPprof(L, 1)
	pp.stop <- true
	return 0
}
