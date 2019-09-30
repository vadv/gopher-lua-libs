# cloudwatch [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/aws/cloudwatch?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/aws/cloudwatch)

## Usage

```lua
local cloudwatch = require("cloudwatch")

local clw_client, err = cloudwatch.new()
if err then error(err) end

local filter = {
    log_group_name = "group-name",
--  filter_patern = "",
    start_time = 1557948000,
    end_time = 1557948200,
}
local timeout_sec = 100
local err = clw_client:download("download.log", filter, timeout_sec)
if err then error(err) end

local query = {
    namespace = "AWS/RDS",
    metric = "CPUUtilization",
    dimension_name = "DBInstanceIdentifier",
    dimension_value = "hostname",
    stat = "Average",
    period = 60,
}
local result, err = clw_client:get_metric_data({start_time=1557948000, end_time=1557948200, queries={cpu=query}})
if err then error(err) end
print(inspect(result))
--[[
{
  cpu = {
    1569880560 = 5.72916666666667,
    1569880620 = 4.29166666666667,
    1569880680 = 4.29583333308498,
    1569880740 = 6.44166666641831,
    1569880800 = 9.30833333358169,
    1569880860 = 5.72500000024835,
    1569880920 = 4.29583333308498,
    1569880980 = 4.29583333308498,
    1569881040 = 5.72500000024835
  },
}
--]]
```
