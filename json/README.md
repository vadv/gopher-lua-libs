# json [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/json?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/json)

## Usage

```lua
local json = require("json")
local inspect = require("inspect")

-- json.encode()
local jsonString = [[
    {
        "a": {"b":1}
    }
]]
local result, err = json.decode(jsonString)
if err then error(err) end
local result = inspect(result, {newline="", indent=""})
if not(result == "{a = {b = 1}}") then error("json.encode") end

-- json.decode()
local table = {a={b=1}}
local result, err = json.encode(table)
if err then error(err) end
local result = inspect(result, {newline="", indent=""})
if not(result == [[{"a":{"b":1}}]]) then error("json.decode") end
```
