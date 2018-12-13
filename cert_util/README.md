# cert_util [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/cert_util?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/cert_util)

## Usage

```lua
local cert_util = require("cert_util")

local tx, err = cert_util.not_after("google.com", "64.233.165.101:443")
if err then error(err) end
if not(tx == 1548838740) then error("error!") end
```
