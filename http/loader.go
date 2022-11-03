// Package http implements golang package http functionality for lua.

package http

import (
	client "github.com/vadv/gopher-lua-libs/http/client"
	server "github.com/vadv/gopher-lua-libs/http/server"
	util "github.com/vadv/gopher-lua-libs/http/util"
	lua "github.com/yuin/gopher-lua"
)

// Preload adds http to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local http = require("http")
func Preload(L *lua.LState) {
	L.PreloadModule("http", Loader)
	client.Preload(L)
	server.Preload(L)
	util.Preload(L)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	http_client_ud := L.NewTypeMetatable(`http_client_ud`)
	L.SetGlobal(`http_client_ud`, http_client_ud)
	L.SetField(http_client_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"do_request": client.DoRequest,
	}))

	http_request_ud := L.NewTypeMetatable(`http_request_ud`)
	L.SetGlobal(`http_request_ud`, http_request_ud)
	L.SetField(http_request_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"set_basic_auth": client.SetBasicAuth,
		"header_set":     client.HeaderSet,
	}))

	http_server_response_writer_ud := L.NewTypeMetatable(`http_server_response_writer_ud`)
	L.SetGlobal(`http_server_response_writer_ud`, http_server_response_writer_ud)
	L.SetField(http_server_response_writer_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"code":     server.HeaderCode,
		"header":   server.Header,
		"write":    server.Write,
		"redirect": server.Redirect,
		"done":     server.Done,
	}))

	http_server_ud := L.NewTypeMetatable(`http_server_ud`)
	L.SetGlobal(`http_server_ud`, http_server_ud)
	L.SetField(http_server_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"accept":             server.Accept,
		"addr":               server.Addr,
		"do_handle_file":     server.HandleFile,
		"do_handle_string":   server.HandleString,
		"do_handle_function": server.HandleFunction,
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"server":         server.New,
	"serve_static":   server.ServeStaticFiles,
	"client":         client.New,
	"request":        client.NewRequest,
	"file_request":   client.NewFileRequest,
	"query_escape":   util.QueryEscape,
	"query_unescape": util.QueryUnescape,
	"parse_url":      util.ParseURL,
	"build_url":      util.BuildURL,
}
