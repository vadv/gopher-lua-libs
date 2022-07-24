local stats = require("stats")

function Test_median(t)
    local result, err = stats.median({ 0, 0, 10 })
    assert(not err, err)
    assert(result == 0, "median get: " .. tostring(result))
end

function Test_percentile(t)
    tests = {
        {
            name = "100 gets 10",
            data = { 0, 0, 10 },
            percentile = 100,
            expected = 10,
        },
        {
            name = "60 gets 0",
            data = { 0, 0, 10 },
            percentile = 60,
            expected = 0,
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            local result, err = stats.percentile(tt.data, tt.percentile)
            assert(not err, err)
            assert(result == tt.expected, "percentile get: " .. tostring(result))
        end)
    end
end

function Test_standard_deviation(t)
    local result, err = stats.standard_deviation({ 1, 1, 1, 1 })
    assert(not err, err)
    assert(result == 0.5, "standard_deviation get: " .. tostring(result))
end
