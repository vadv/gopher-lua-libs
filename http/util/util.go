package util

import (
	"net/url"

	lua "github.com/yuin/gopher-lua"
)

// QueryEscape(): lua http.query_escape(string) returns escaped string
func QueryEscape(L *lua.LState) int {
	query := L.CheckString(1)
	escapedUrl := url.QueryEscape(query)
	L.Push(lua.LString(escapedUrl))
	return 1
}

// QueryUnescape(): lua http.query_unescape(string) returns unescaped (string, error)
func QueryUnescape(L *lua.LState) int {
	query := L.CheckString(1)
	url, err := url.QueryUnescape(query)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(url))
	return 1
}
