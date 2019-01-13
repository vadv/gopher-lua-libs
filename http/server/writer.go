package http

import (
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type luaWriter struct {
	http.ResponseWriter
	req      *http.Request
	complete bool
	done     chan bool
}

// NewWriter return lua userdata with luaWriter
func NewWriter(L *lua.LState, w http.ResponseWriter, req *http.Request, done chan bool) *lua.LUserData {
	luaWriter := &luaWriter{ResponseWriter: w, done: done, req: req}
	ud := L.NewUserData()
	ud.Value = luaWriter
	L.SetMetatable(ud, L.GetTypeMetatable("http_server_response_writer_ud"))
	return ud
}

func checkServeWriter(L *lua.LState, n int) *luaWriter {
	ud := L.CheckUserData(n)
	w, ok := ud.Value.(*luaWriter)
	if !ok {
		L.ArgError(1, "must be http_server_response_writer_ud")
	}
	return w
}

// HeaderCode lua http_server_response_writer_ud:code(number)
func HeaderCode(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	code := L.CheckInt(2)
	w.ResponseWriter.WriteHeader(code)
	return 0
}

// Header lua http_server_response_writer_ud:header(key, value)
func Header(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	key, value := L.CheckString(2), L.CheckString(3)
	w.Header().Set(key, value)
	return 0
}

// Write lua http_server_response_writer_ud:write(string) return (number, err)
func Write(L *lua.LState) int {
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

// Redirect lua http_server_response_writer_ud:redirect(url, code) return err
func Redirect(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	url := L.CheckString(2)
	code := http.StatusPermanentRedirect
	if L.GetTop() > 2 {
		code = L.CheckInt(3)
	}
	http.Redirect(w.ResponseWriter, w.req, url, code)
	return 0
}

// Done lua http_server_response_writer_ud:done()
func Done(L *lua.LState) int {
	w := checkServeWriter(L, 1)
	w.complete = true
	w.done <- true
	return 0
}
