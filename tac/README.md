# tac [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/tac?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/tac)

## Usage

```lua
local file = io.open("data.txt", "w")
file:write("1", "\n")
file:write("2", "\n")
file:write("3", "\n")

local tac = require("tac")
local scanner, err = tac.open("data.txt")
if err then error(err) end

while true do
    local line = scanner:line()
    if line == nil then break end
    print(line)
end
scanner:close()

-- Output:
-- 3
-- 2
-- 1
```

