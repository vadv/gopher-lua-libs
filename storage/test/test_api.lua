local storage = require("storage")
local inspect = require("inspect")
local time = require("time")

function Test_storage(t)
    local s, err = storage.open("./test/db.json")
    assert(not err, err)

    t:Run("set", function(t)
        local err = s:set("key", { "one", "two", 1 }, 1)
        assert(not err, err)

        local err = s:set("key2", "value2", 60)
        assert(not err, err)

        local err = s:set("key3", 10.64, nil)
        assert(not err, err)

        local value, found, err = s:get("key")
        assert(not err, err)
        assert(found, "must be found")

        assert(value[1] == "one", "value: " .. value[1])
        assert(value[3] == 1, "value: " .. value[3])

        local value, found, err = s:get("key2")
        assert(not err, err)
        assert(value == "value2", "value: " .. value)

        local value, found, err = s:get("key3")
        assert(not err, err)
        assert(value == 10.64, "value: " .. value)
    end)

    t:Run("after ttl", function(t)
        time.sleep(1)

        -- wait ttl
        local value, found, err = s:get("key")
        assert(not err, err)
        assert(not found, "must be not found")
    end)

    t:Run("close", function(t)
        -- close
        local err = s:close()
        assert(not err, err)
    end)

    t:Run("get nil", function(t)
        local err = s:set("key2", nil, 60)
        assert(not err, err)
        local value, found, err = s:get("key")
        assert(not err, err)
        assert(not found, "must be not found")
        assert(value == nil, "value must be nil, but was: " .. tostring(value))
    end)

    t:Run("keys", function(t)
        local keys = s:keys()
        assert(#keys == 2, "keys: " .. #keys)
    end)

    t:Run("dump", function(t)
        local dump, err = s:dump()
        assert(not err, err)
        assert(dump.key3 == 10.64, "dump: " .. tostring(dump.key3))
    end)
end