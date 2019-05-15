# cloudwatch [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/aws/cloudwatch?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/aws/cloudwatch)

## Usage

```lua
local cloudwatch = require("cloudwatch")

local clw_client, err = cloudwatch.new()
if err then error(err) end

local filter = {
    log_group_name = "group-name",
    start_time = 1557948000,
    end_time = 1557948200,
}
local timeout_sec = 100
local err = clw_client:download("download.log", filter, timeout_sec)
if err then error(err) end
```
