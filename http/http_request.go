package http

import (
	"bytes"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

type luaRequest struct {
	*http.Request
}

func checkRequest(L *lua.LState, n int) *luaRequest {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaRequest); ok {
		return v
	}
	L.ArgError(1, "http request excepted")
	return nil
}

// http.request(verb, url, body) returns user-data, error
func NewRequest(L *lua.LState) int {
	verb := L.CheckString(1)
	url := L.CheckString(2)
	buffer := &bytes.Buffer{}
	if L.GetTop() > 2 {
		buffer.WriteString(L.CheckString(3))
	}
	httpReq, err := http.NewRequest(verb, url, buffer)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	req := &luaRequest{Request: httpReq}
	req.Request.Header.Set(`User-Agent`, DefaultUserAgent)
	ud := L.NewUserData()
	ud.Value = req
	L.SetMetatable(ud, L.GetTypeMetatable("http_request_ud"))
	L.Push(ud)
	return 1
}

// request:set_basic_auth(username, password)
func RequestSetBasicAuth(L *lua.LState) int {
	req := checkRequest(L, 1)
	user, passwd := L.CheckAny(2).String(), L.CheckAny(3).String()
	req.SetBasicAuth(user, passwd)
	return 0
}

// request:header_set(key, value)
func RequestHeaderSet(L *lua.LState) int {
	req := checkRequest(L, 1)
	key, value := L.CheckAny(2).String(), L.CheckAny(3).String()
	req.Header.Set(key, value)
	return 0
}
