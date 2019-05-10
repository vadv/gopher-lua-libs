# cmd [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/cmd?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/cmd)

## Usage

```lua
local cmd = require("cmd")
local runtime = require("runtime")

local command = "sleep 1"
if runtime.goos() == "windows" then command = "timeout 1" end

local result, err = cmd.exec(command)
if err then error(err) end
if not(result.status == 0) then error("status") end
```

The default timeout is 10 seconds after which the command will be terminated. The default timeout may be overriden with an optional timeout value (seconds).

```lua
local cmd = require("cmd")
local runtime = require("runtime")

local command = "sleep 5"
if runtime.goos() == "windows" then command = "timeout 1" end

local result, err = cmd.exec(command, 1)
if err then error(err) end
if not(result.status == 0) then error("status") end
```
