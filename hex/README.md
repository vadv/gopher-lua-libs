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

s, err = hex.decode_string("Zm9vAWJhcg")
assert(not err, err)
print(s)
foobar

s, err = hex.StdEncoding:decode_string("Zm9vAWJhcg==")
assert(not err, err)
print(s)
foobar

s, err = hex.RawURLEncoding:decode_string("dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA")
assert(not err, err)
print(s)
this is a <tag> and should be encoded

s, err = hex.URLEncoding:decode_string("dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==")
assert(not err, err)
print(s)
this is a <tag> and should be encoded
```
