# pb [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/strings?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/pb)

## Usage

```lua
local pb = require('pb')
local time = require('time')

local count = 100
local bar = pb.new(count)

bar:start()

for i=1, count, 1 do
  time.sleep(1)
  bar:increment()
end

bar:finish()
```

### Configure progress bar
- `bar:configure({})` - change progress bar parameters. Avaliable options:
```
template - use custom template (Please see https://github.com/cheggaaa/pb/blob/master/v3/element.go for all available elements)
refresh_rate in ms (default 200ms)
writer (default stderr). Only supported value is `stdout`.
```
