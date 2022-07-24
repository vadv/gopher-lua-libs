# inspect [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/inspect?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/inspect)

## Usage

```lua
-- script.lua
local argparse = require "argparse"

local parser = argparse("script", "An example.")
parser:argument("input", "Input file.")
parser:option("-o --output", "Output file.", "a.out")
parser:option("-I --include", "Include locations."):count("*")

local args = parser:parse()
```
