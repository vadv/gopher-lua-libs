local log = require("log")
local ioutil = require("ioutil")

os.remove("./test/test.log")

local info, err = log.new("./test/test.log")
if err then error(err) end

info:println("1", 2)

-- check set prefix
info:set_prefix("[INFO] ")
info:printf("%s %.2f", "1", 2.2)

-- check flags
info:set_flags( {longfile=true} )
info:printf("ok\n")
info:print("ok")
info:println("ok")

local err = info:close()
if err then error(err) end

-- check result

local except_result = [[1 2
[INFO] 1 2.20
[INFO] ./test/test_api.lua:17: ok
[INFO] ./test/test_api.lua:18: ok
[INFO] ./test/test_api.lua:19: ok
]]

local get_result, err = ioutil.read_file("./test/test.log")
if err then error(err) end

if not(except_result == get_result) then
    error("log fail")
end
