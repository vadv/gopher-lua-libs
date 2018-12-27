# goos [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/goos?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/goos)

## Usage

```lua
local goos = require("goos")

local stat, err = goos.stat("./filename")
if err then error(err) end

print(stat.is_dir)
print(stat.size)
print(stat.mod_time)
print(stat.mode)
```

