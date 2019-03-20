package pprof_test

import (
	"log"

	lua_http "github.com/vadv/gopher-lua-libs/http"
	lua_pprof "github.com/vadv/gopher-lua-libs/pprof"
	lua_time "github.com/vadv/gopher-lua-libs/time"

	lua "github.com/yuin/gopher-lua"
)

// pprof:register(), pprof_ud:enable(), pprof_ud:disable()
func Example_package() {
	state := lua.NewState()
	lua_pprof.Preload(state)
	lua_http.Preload(state)
	lua_time.Preload(state)
	source := `
local pprof = require("pprof")
local http = require("http")
local time = require("time")

local client = http.client()
local pp = pprof.register(":1234")

pp:enable()
time.sleep(1)

local req, err = http.request("GET", "http://127.0.0.1:1234/debug/pprof/goroutine")
if err then error(err) end
local resp, err = client:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
print(resp.code)

pp:disable()
time.sleep(5)

local resp, err = client:do_request(req)
if not(err) then error("must be error") end
        `
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 200
}
