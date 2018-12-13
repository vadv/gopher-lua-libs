// Package http implements golang package http functionality for lua.
package http

import (
	"io/ioutil"
	"net/url"

	lua "github.com/yuin/gopher-lua"
)

// DoRequest(): lua http_client_ud:do_request()
// http_client_ud:do_request(http_request_ud) returns (response, error)
//    response: {
//      code=http_code (200, 201, ..., 500, ...),
//      body=string
//    }
func DoRequest(L *lua.LState) int {
	client := checkClient(L)
	req := checkRequest(L, 2)
	response, err := client.Do(req.Request)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.NewTable()
	L.SetField(result, `code`, lua.LNumber(response.StatusCode))
	L.SetField(result, `body`, lua.LString(string(data)))
	L.Push(result)
	return 1
}

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
