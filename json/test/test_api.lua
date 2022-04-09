local json = require("json")
local inspect = require("inspect")

function TestJson(t)
    local jsonStringWithNull = [[{"a":{"b":1, "c":null}}]]
    local jsonString = [[{"a":{"b":1}}]]

    local result, err = json.decode(jsonStringWithNull)
    t:Run("decode", function(t)
        if err then error(err) end
        if result["a"]["b"] ~= 1 then error("must be decode") end
        if result["a"]["c"] ~= nil then error("c is not nil") end
        print("done: json.decode()")
    end)

    local result, err = json.encode(result)
    t:Run("encode omits null values", function(t)
        if err then error(err) end
        if result ~= jsonString then error("must be encode "..inspect(result)) end
        print("done: json.encode()")
    end)
end