# hex [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/hex?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/hex)

Lua module for [encoding/hex](https://pkg.go.dev/encoding/hex)

## Usage

### Encoding

```lua
local hex = require("hex")

s = hex:encode_to_string("foo\01bar")
print(s)
666f6f01626172

```

### Decoding

```lua
local hex = require 'hex'

decoded, err = hex.decode_string("666f6f01626172")
assert(not err, err)
print(decoded)
foobar

encoded = hex.encode_to_string(decoded)
print(encoded)
666f6f01626172
```
