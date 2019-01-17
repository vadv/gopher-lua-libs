package http

import (
	"io/ioutil"
	"net/http"

	lua "github.com/yuin/gopher-lua"
)

// NewRequest return lua table with http.Request representation
func NewRequest(L *lua.LState, req *http.Request) *lua.LTable {
	luaRequest := L.NewTable()
	bodyReader := L.NewFunction(func(L *lua.LState) int {
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(lua.LString(string(data)))
		return 1
	})
	luaRequest.RawSetString(`body`, bodyReader)
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
