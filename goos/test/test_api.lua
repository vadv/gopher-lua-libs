local goos = require("goos")
local runtime = require("runtime")

local info, err = goos.stat("./test/test.file")
if err then error(err) end

if not(info.is_dir == false) then error("is dir") end
if not(0 == info.size) then error("size") end
if not(info.mod_time > 0) then error("mod_time") end

if runtime.goos() == "linux" then
  if not("-rwxrwxrwx" == info.mode) then error("mode") end
else
  if not("-rw-rw-rw-" == info.mode) then error("mode") end
end
