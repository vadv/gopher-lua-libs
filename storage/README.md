# storage [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/storage?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/storage)

## Usage

```lua
local storage = require("storage")

local s, err = storage.open("./test/db.json")
if err then error(err) end

-- key, value, ttl (default = 60s)
local err = s:set("key", {"one", "two", 1}, 10)
if err then error(err) end

local value, found, err = s:get("key")
if err then error(err) end
if not found then error("must be found") end

if not(value[1] == "one") then error("value") end
if not(value[3] == 1) then error("value") end

-- override with set max ttl
local err = s:set("key", "override", nil)
local value, found, err = s:get("key")
if not(value == "ovveride") then error("must be found") end
```

