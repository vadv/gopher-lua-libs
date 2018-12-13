# tcp [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/tcp?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/tcp)

## Usage

```lua
local tcp = require("tcp")

-- http request
local conn, err = tcp.open("google.com:80")
err = conn:write("GET /\n\n")
if err then error(err) end
local result, err = conn:read(64*1024)
print(result)

-- ping pong game
local conn, err = tcp.open(":12345")
if err then error(err) end

err = conn:write("ping")
if err then error(err) end

local result, err = conn:read()
if err then error(err) end
if (result == "pong") then error("must be pong message") end
```

