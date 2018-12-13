local filepath = require("filepath")

local path = "1"
local need_path = path .. filepath.separator() .. "2" .. filepath.separator() .. "3"
path = filepath.join(path, "2", "3")
if not(path == need_path) then error("filepath.join()") end
print("done: filepath.join(), filepath.separator()")

local t = filepath.glob("test"..filepath.separator().."*")
local count_t = 0
for k, v in pairs( t ) do
    count_t = count_t + 1
    if k == 1 then if not(v == "test"..filepath.separator().."test_api.lua") then error("filepath.glob()") end end
end
if not(count_t == 1) then error("filepath.join()") end
print("done: filepath.glob()")
