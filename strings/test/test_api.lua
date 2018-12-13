local strings = require("strings")

local str = "hello world"

local t = strings.split(str, " ")
local count_t = 0
for k, v in pairs(t) do
    count_t = count_t + 1
    if k == 1 then if not(v == "hello") then error("strings.split()") end end
    if k == 2 then if not(v == "world") then error("strings.split()") end end
end
if not(count_t == 2) then error("string.split()") end
print("done: strings.split()")

if not(strings.has_prefix(str, "hello")) then error("strings.has_prefix()") end
if not(strings.has_suffix(str, "world")) then error("strings.has_suffix()") end
print("done: strings.has_suffix, strings.has_prefix")

if not(strings.trim(str, "world") == "hello ") then error("strings.trim()") end
if not(strings.trim(str, "hello ") == "world") then error("strings.trim()") end
print("done: strings.trim()")
