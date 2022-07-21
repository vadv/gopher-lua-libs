if os.getenv("CI") then
    -- travis: Include a test function that just skips
    function TestCI(t)
        t:Skip("CI")
    end
else
    --dofile("./test/test_cloudwatch_logs.lua")
    dofile("./test/test_cloudwatch_get_metric_data.lua")
end
