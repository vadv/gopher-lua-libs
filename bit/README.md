# stats [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/bit?bit.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/bit)

## Usage

```lua
local bit = require("bit")

local result, _ = bit.band(1, 0)
print(result)
-- Output: 0

local result, _ = bit.lshift(10, 5)
print(result)
-- Output: 320
```

