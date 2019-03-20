# pprof [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/pprof?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/pprof)

## Usage

```lua
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

pp:disable()
time.sleep(5)

local resp, err = client:do_request(req)
if not(err) then error("must be error") end
```
