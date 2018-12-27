local goos = require("goos")

-- stat
local info, err = goos.stat("./test/test.file")
if err then error(err) end
if not(info.is_dir == false) then error("is dir") end
if not(0 == info.size) then error("size") end
if not(info.mod_time > 0) then error("mod_time") end
if not(info.mode > "") then error("mode") end

-- hostname
local hostname, err = goos.hostname()
if err then error(err) end
if not(hostname > "") then error("hostname") end

-- get_pagesize
if not(goos.get_pagesize() > 0) then error("pagesize") end
