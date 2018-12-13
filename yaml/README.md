# yaml [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/yaml?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/yaml)

## Usage

```lua
local yaml = require("yaml")
local inspect = require("inspect")

-- yaml.decode()
local text = [[
a:
  b: 1
]]
local result, err = yaml.decode(text)
if err then error(err) end
print(inspect(result, {newline="", indent=""}))
-- Output:
-- {a = {b = 1}}
```

