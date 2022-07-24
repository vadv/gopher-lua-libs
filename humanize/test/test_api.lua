local humanize = require("humanize")
local time = require("time")

function Test_parse_bytes(t)
    local size, err = humanize.parse_bytes("1.3GiB")
    assert(not err, err)
    assert(size == 1395864371, "size: " .. tostring(size))
end

function Test_ibytes(t)
    local size_string = humanize.ibytes(1395864371)
    assert(size_string == "1.3 GiB", "ibytes: " .. size_string)
end

function Test_time(t)
    local t = time.unix() - 2
    local time_string = humanize.time(t)
    assert(time_string == "2 seconds ago", "time: " .. time_string)
end

function Test_si(t)
    local si_result = humanize.si(1202121, "m")
    assert(si_result == "1.202121 Mm", "si: " .. tostring(si_result))
end
