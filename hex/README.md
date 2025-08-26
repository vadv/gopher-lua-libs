# hex [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/hex?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/hex)

Lua module for [encoding/hex](https://pkg.go.dev/encoding/hex)

## Usage

### Encoding

```lua
local hex = require("hex")

s = hex.RawStdEncoding:encode_to_string("foo\01bar")
print(s)
Zm9vAWJhcg

s = hex.StdEncoding:encode_to_string("foo\01bar")
print(s)
Zm9vAWJhcg==

s = hex.RawURLEncoding:encode_to_string("this is a <tag> and should be encoded")
print(s)
dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA

s = hex.URLEncoding:encode_to_string("this is a <tag> and should be encoded")
print(s)
dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==

```

### Decoding

```lua
local hex = require 'hex'

decoded, err = hex.decode_string("666f6f62617262617a")
assert(not err, err)
print(decoded)
foobar

encoded = hex.encode_to_string(decoded)
print(encoded)
666f6f62617262617a
```
