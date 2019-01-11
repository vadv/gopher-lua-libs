package http

import (
	"net"
	"net/http"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type luaServer struct {
	*http.Server
	net.Listener
	sync.Mutex
	serveData chan *serveData
	err       error
}

type serveData struct {
	w    http.ResponseWriter
	req  *http.Request
	done chan bool
}

type luaServeWriter struct {
	http.ResponseWriter
	req  *http.Request
	done chan bool
}

func checkServeWriter(L *lua.LState, n int) *luaServeWriter {
	ud := L.CheckUserData(n)
	w, ok := ud.Value.(*luaServeWriter)
	if !ok {
		L.ArgError(1, "must be http_server_response_writer_ud")
	}
	return w
}

// serveWriteHeaderCode lua http_server_response_writer_ud:code(number)
func serveWriteHeaderCode(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	code := L.CheckInt(2)
	w.ResponseWriter.WriteHeader(code)
	return 0
}

// serveWriteHeader lua http_server_response_writer_ud:header(key, value)
func serveWriteHeader(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	key, value := L.CheckString(2), L.CheckString(3)
	w.Header().Set(key, value)
	return 0
}

// serveWrite lua http_server_response_writer_ud:write(string) return (number, err)
func serveWrite(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	data := L.CheckAny(2).String()
	count, err := w.ResponseWriter.Write([]byte(data))
	L.Push(lua.LNumber(count))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// serveRedirect lua http_server_response_writer_ud:redirect(url, code) return err
func serveRedirect(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	url := L.CheckString(2)
	code := http.StatusPermanentRedirect
	if L.GetTop() > 2 {
		code = L.CheckInt(3)
	}
	http.Redirect(w.ResponseWriter, w.req, url, code)
	return 0
}

// serveDone lua http_server_response_writer_ud:done()
func serveDone(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	w.done <- true
	return 0
}

func checkServer(L *lua.LState, n int) *luaServer {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaServer); ok {
		return v
	}
	L.ArgError(n, "http server excepted")
	return nil
}

// run serve
func (s *luaServer) serve() {
	s.err = http.Serve(s.Listener, s)
}

// http.server(bind, handler) returns (user data, error)
func NewServer(L *lua.LState) int {
	bind := L.CheckAny(1).String()
	l, err := net.Listen(`tcp`, bind)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	server := &luaServer{
		Listener:  l,
		serveData: make(chan *serveData, 1),
	}
	go server.serve()
	ud := L.NewUserData()
	ud.Value = server
	L.SetMetatable(ud, L.GetTypeMetatable("http_server_ud"))
	L.Push(ud)
	return 1
}

// NewLuaRequest return lua table with http.Request representation
func NewLuaRequest(L *lua.LState, req *http.Request) *lua.LTable {
	luaRequest := L.NewTable()
	luaRequest.RawSetString(`host`, lua.LString(req.Host))
	luaRequest.RawSetString(`method`, lua.LString(req.Method))
	luaRequest.RawSetString(`referer`, lua.LString(req.Referer()))
	luaRequest.RawSetString(`proto`, lua.LString(req.Proto))
	luaRequest.RawSetString(`user_agent`, lua.LString(req.UserAgent()))
	if req.URL != nil && len(req.URL.Query()) > 0 {
		query := L.NewTable()
		for k, v := range req.URL.Query() {
			if len(v) > 0 {
				query.RawSetString(k, lua.LString(v[0]))
			}
		}
		luaRequest.RawSetString(`query`, query)
	}
	if len(req.Header) > 0 {
		headers := L.NewTable()
		for k, v := range req.Header {
			if len(v) > 0 {
				headers.RawSetString(k, lua.LString(v[0]))
			}
		}
		luaRequest.RawSetString(`headers`, headers)
	}
	luaRequest.RawSetString(`path`, lua.LString(req.URL.Path))
	luaRequest.RawSetString(`raw_path`, lua.LString(req.URL.RawPath))
	luaRequest.RawSetString(`raw_query`, lua.LString(req.URL.RawQuery))
	luaRequest.RawSetString(`request_uri`, lua.LString(req.RequestURI))
	luaRequest.RawSetString(`remote_addr`, lua.LString(req.RemoteAddr))
	return luaRequest
}

// ServerAccept lua http_server_ud:accept() returns request_table, http_server_response_writer_ud
func ServerAccept(L *lua.LState) int {
	s := checkServer(L, 1)
	select {
	case data := <-s.serveData:

		// make request
		luaRequest := NewLuaRequest(L, data.req)

		// make writer
		luaWriter := &luaServeWriter{ResponseWriter: data.w, done: data.done, req: data.req}
		ud := L.NewUserData()
		ud.Value = luaWriter

		L.SetMetatable(ud, L.GetTypeMetatable("http_server_response_writer_ud"))
		L.Push(luaRequest)
		L.Push(ud)
		return 2
	}
}

// ServeHTTP interface realisation
func (s *luaServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	doneChan := make(chan bool)
	data := &serveData{w: w, req: req, done: doneChan}
	// send data for lua
	s.serveData <- data
	// wait response from lua
	select {
	case <-doneChan:
		return
	}
}
