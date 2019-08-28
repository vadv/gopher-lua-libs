local stats = require("stats")

local result, err = stats.median({0,0,10})
if err then error(err) end
if not(result == 0) then error("median get: "..tostring(result)) end

local result, err = stats.percentile({0,0,10}, 100)
if err then error(err) end
if not(result == 10) then error("percentile get: "..tostring(result)) end

local result, err = stats.percentile({0,0,10}, 60)
if err then error(err) end
if not(result == 0) then error("percentile get: "..tostring(result)) end

local result, err = stats.standard_deviation({1,1,1,1})
if err then error(err) end
if not(result == 0.5) then error("standard_deviation get: "..tostring(result)) end