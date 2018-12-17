local json = require("json")
local inspect = require("inspect")

local jsonStringWithNull = [[{"a":{"b":1, "c":null}}]]
local jsonString = [[{"a":{"b":1}}]]

local result, err = json.decode(jsonStringWithNull)
if err then error(err) end
if not(result["a"]["b"] == 1) then error("must be decode") end
print("done: json.decode()")

local result, err = json.encode(result)
if err then error(err) end
if not(result==jsonString) then error("must be encode "..inspect(result)) end
print("done: json.encode()")
