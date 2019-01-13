package http_test

import (
	"log"

	http "github.com/vadv/gopher-lua-libs/http"
	plugin "github.com/vadv/gopher-lua-libs/plugin"
	lua "github.com/yuin/gopher-lua"
)

// http:server()
func ExampleAccept() {
	state := lua.NewState()
	http.Preload(state)
	plugin.Preload(state)
	source := `
    local http = require("http")
    local plugin = require("plugin")

    local server, err = http.server("127.0.0.1:1999")
    if err then error(err) end

    local client_plugin = [[
        local time = require("time")
        local http = require("http")
        time.sleep(1)
        local client = http.client({timeout=1})
        local request, err = http.request("GET", "http://127.0.0.1:1999/get/url?param1=value1")
        if err then error(err) end
        client:do_request(request)
    ]]

    local client_plugin = plugin.do_string(client_plugin)
    client_plugin:run()

    local req, resp = server:accept()
    print("host: "..req.host)
    print("method: "..req.method)
    -- print("referer: "..req.referer)
    print("proto: "..req.proto)
    print("request_uri: "..req.request_uri)
    print("user_agent: "..req.user_agent)
    -- print(req.remote_addr)
    print("header: Accept-Encoding="..req.headers["Accept-Encoding"])
    for k, v in pairs(req.query) do
      print("query: "..k.."="..v)
    end
    resp:code(200) -- write code
    resp:header("Content-Type", "application/json") -- write header
    resp:write("ok")
    resp:done()

    client_plugin:stop()
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// host: 127.0.0.1:1999
	// method: GET
	// proto: HTTP/1.1
	// request_uri: /get/url?param1=value1
	// user_agent: gopher-lua
	// header: Accept-Encoding=gzip
	// query: param1=value1
}
