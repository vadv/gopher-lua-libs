local json = require("json")

local jsonString = [[{"a":{"b":1}}]]

local result, err = json.encode(jsonString)
if err then error(err) end
if not(result["a"]["b"] == 1) then error("must be encode") end
print("done: json.encode()")

local result, err = json.decode(result)
if err then error(err) end
if not(result==jsonString) then error("must be decode") end
print("done: json.decode()")
