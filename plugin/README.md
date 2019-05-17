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

-- plugin.do_string(body)
-- also you can use: plugin.do_file(filename)
local print_plugin = plugin.do_string(plugin_body)
print_plugin:run()
time.sleep(2)
print_plugin:stop()
time.sleep(1)

local running = print_plugin:is_running()
if running then error("already running") end
-- also you can get last error: print_plugin:error()

-- plugin.do_string_with_payload()
-- also you can use: plugin.do_file_with_payload(filename)
local job_body = [[
    print(payload)
]]
local print_plugin_with_payload = plugin.do_string_with_payload(plugin_body, "text of payload")
print_plugin_with_payload:run()
-- must print: "text of payload"
time.sleep(1)
local running = print_plugin:is_running()
if running then error("already running") end
```
