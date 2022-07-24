function Test_cloudwatch_logs(t)
    local cloudwatch = require("cloudwatch")

    local clw_client, err = cloudwatch.new()
    assert(not err, err)

    local filter = {
        log_group_name = os.getenv("LOG_GROUP_NAME"),
        start_time = 1557948000,
        end_time = 1557948200,
    }
    local err = clw_client:download("./test/test.log", filter, 100)
    assert(not err, err)
end
