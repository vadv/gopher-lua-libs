local humanize = require("humanize")
local time = require("time")

-- parse_bytes
local size, err = humanize.parse_bytes("1.3GiB")
if err then error(err) end
if not(size == 1395864371) then error("size: "..tostring(size)) end

-- ibytes
local size_string = humanize.ibytes(1395864371)
if not(size_string == "1.3 GiB") then error("ibytes: "..size_string) end

-- time
local t = time.unix() - 2
local time_string = humanize.time(t)
if not(time_string == "2 seconds ago") then error("time: "..time_string) end
