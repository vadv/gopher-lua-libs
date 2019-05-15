# ioutil [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/ioutil?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/ioutil)

## Usage

```lua
local ioutil = require("ioutil")

-- ioutil.write_file()
local err = ioutil.write_file("./test/file.data", "content of test file")
if err then error(err) end

-- ioutil.read_file()
local result, err = ioutil.read_file("./test/file.data")
if err then error(err) end
if not(result == "content of test file") then error("ioutil.read_file()") end
```

