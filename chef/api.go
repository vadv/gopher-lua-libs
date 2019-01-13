// Package chef implements chef client api functionality for lua.
package chef

import (
	"bytes"
	"fmt"
	"net/url"

	lua_http "github.com/vadv/gopher-lua-libs/http/client/interface"
	lua_json "github.com/vadv/gopher-lua-libs/json"

	lua "github.com/yuin/gopher-lua"
)

// NewClient lua chef.client(client_name, path_to_file_with_key, chef_url, http_client_ud) returns (chef_client_ud, err)
func NewClient(L *lua.LState) int {
	name := L.CheckString(1)
	filename := L.CheckString(2)
	urlStr := L.CheckString(3)
	chefClient := &luaChefClient{name: name}
	pk, err := privateKeyFromFile(filename)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	url, err := url.Parse(urlStr)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	chefClient.url = url
	chefClient.key = pk
	if L.GetTop() > 3 {
		// http client
		ud := L.CheckUserData(4)
		if v, ok := ud.Value.(lua_http.LuaHTTPClient); ok {
			chefClient.LuaHTTPClient = v
		} else {
			L.ArgError(2, "must be http_client_ud")
		}
	} else {
		chefClient.LuaHTTPClient = lua_http.NewPureClient()
	}
	ud := L.NewUserData()
	ud.Value = chefClient
	L.SetMetatable(ud, L.GetTypeMetatable(`chef_client_ud`))
	L.Push(ud)
	return 1
}

// Request lua chef_client_ud:request("GET|POST|PUT", "/api/path", "body") returns (table, error)
func Request(L *lua.LState) int {
	client := checkChefClient(L, 1)
	verb := L.CheckString(2)
	url := L.CheckString(3)
	var body []byte
	if L.GetTop() > 3 {
		data := L.CheckAny(4).String()
		body = []byte(data)
	}
	responseBytes, err := client.request(verb, url, bytes.NewReader(body))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	value, err := lua_json.ValueDecode(L, responseBytes)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(value)
	return 1
}

// Search lua chef_client_ud:search(
// "index",
// "query",
// "partical"={
//   "return_name" = ["node_attribute_name"],
// },
// params={
//    sort_by = "X_CHEF_id_CHEF_X asc",
//    start = 0,
//    rows  = 1000
//}
//) returns (table, error)
func Search(L *lua.LState) int {
	client := checkChefClient(L, 1)
	index := L.CheckString(2)
	query := L.CheckString(3)
	// partical data
	var partical *lua.LTable
	var particalJSON []byte
	if L.GetTop() > 3 {
		value := L.CheckAny(4)
		if value != lua.LNil {
			partical = L.CheckTable(4)
		}
	}
	if partical != nil {
		data, err := lua_json.ValueEncode(partical)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		particalJSON = data
	}
	// parse params
	param := L.NewTable()
	if L.GetTop() > 4 {
		param = L.CheckTable(5)
	}
	// default values in param
	start, rows, sortBy := 0, 1000, "X_CHEF_id_CHEF_X asc"
	param.ForEach(func(k lua.LValue, v lua.LValue) {
		// parse start
		if k.String() == `start` {
			if value, ok := v.(lua.LNumber); ok {
				start = int(value)
			} else {
				L.ArgError(4, "start must be number")
			}
		}
		// parse rows
		if k.String() == `rows` {
			if value, ok := v.(lua.LNumber); ok {
				rows = int(value)
			} else {
				L.ArgError(4, "rows must be number")
			}
		}
		// parse sort_by
		if k.String() == `sort_by` {
			if _, ok := v.(lua.LString); ok {
				sortBy = v.String()
			} else {
				L.ArgError(1, "sort_by must be string")
			}
		}
	})

	// make search url
	queryString := fmt.Sprintf("search/%s?q=%s&rows=%d&sort=%s&start=%d", index, query, rows, sortBy, start)

	verb := "GET"
	if particalJSON != nil {
		// if have body make post request
		verb = "POST"
	}
	responseBytes, err := client.request(verb, queryString, bytes.NewReader(particalJSON))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	value, err := lua_json.ValueDecode(L, responseBytes)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(value)
	return 1
}
