package http

import (
	"net"
	"net/http"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// ServeStaticFiles lua http:serve_static("directory", ":port") return err
func ServeStaticFiles(L *lua.LState) int {
	staticDir := L.CheckString(1)
	addr := L.CheckString(2)
	fs := http.FileServer(http.Dir(staticDir))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	server := &http.Server{Handler: fs, IdleTimeout: time.Second * 60}

	// shutdown
	go func(L *lua.LState, l net.Listener) {
		ctx := L.Context()
		if ctx != nil {
			select {
			case <-ctx.Done():
				l.Close()
			}
		}
	}(L, listener)

	err = server.Serve(listener)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
