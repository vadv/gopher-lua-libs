local time = require("time")

local lua_before = os.clock()
local before = time.unix()
time.sleep(2)
local after = time.unix()
local lua_after = os.clock()
if after - before < 1 then error("time.unix()") end
if lua_after - lua_before < 2 then error("time.sleep()") end
print("done: time.sleep(), time.unix()")

local parse, err = time.parse("Dec  2 03:33:05 2018", "Jan  2 15:04:05 2006")
if err then error(err) end
if not(parse == 1543721585) then error("time.parse(): 1") end
print("done: time.parse(): 1")

local _, err = time.parse("Dec  32 03:33:05 2018", "Jan  2 15:04:05 2006")
if (err == nil) then error("time.parse(): must be error") end
print("done: time.parse(): 2")
