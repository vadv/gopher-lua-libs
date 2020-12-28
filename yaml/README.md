# yaml [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/yaml?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/yaml)

## Usage

### decode
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

### encode
```lua
    local yaml = require("yaml")
    local inspect = require("inspect")
    local encoded, err = yaml.encode({a = {b = 1}})
    if err then error(err) end
    print(encoded)
	-- Output:
	-- a:
	--   b: 1
	--
```
