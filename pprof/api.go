// Package pprof implements golang package pprof functionality.
package pprof

import (
	"context"
	"log"
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

// Create(string): return (pprof_ud, err)
func Create(L *lua.LState) int {
	addr := L.CheckString(1)
	ud := L.NewUserData()
	ud.Value = &luaPprof{addr: addr, stop: make(chan bool, 1)}
	L.SetMetatable(ud, L.GetTypeMetatable(`pprof_ud`))
	L.Push(ud)
	return 1
}

// Start(): start pprof
func Start(L *lua.LState) int {
	pp := checkPprof(L, 1)
	go func() {
		h := &http.Server{Addr: pp.addr}
		log.Printf("[INFO] start pprof server at: %s\n", pp.addr)
		go func() {
			if err := h.ListenAndServe(); err != nil {
				log.Printf("[ERROR] pprof at %s: %s\n", pp.addr, err.Error())
				return
			}
		}()
		<-pp.stop
		log.Printf("[INFO] stop pprof: %s\n", pp.addr)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		h.Shutdown(ctx)
	}()
	return 0
}

// Stop(): pprof stop
func Stop(L *lua.LState) int {
	pp := checkPprof(L, 1)
	pp.stop <- true
	log.Printf("[INFO] send stop to pprof server at: %s\n", pp.addr)
	return 0
}
