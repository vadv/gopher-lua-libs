# plugin [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/plugin?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/plugin)

## Usage

```lua
local plugin = require("plugin")
local time = require("time")

local plugin_body = [[
    local time = require("time")
    local i = 1
    while true do
        print(i)
        i = i + 1
        time.sleep(1)
    end
]]

local print_plugin = plugin.new(plugin_body)
print_plugin:run()
time.sleep(2)
print_plugin:stop()
time.sleep(1)

local running = print_plugin:is_running()
if running then error("already running") end
```

