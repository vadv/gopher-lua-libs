local cloudwatch = require("cloudwatch")

local clw_client, err = cloudwatch.new()
if err then error(err) end

local filter = {
    log_group_name = os.getenv("LOG_GROUP_NAME"),
    start_time = 1557948000,
    end_time = 1557948200,
}
local err = clw_client:download("./test/test.log", filter, 100)
if err then error(err) end
