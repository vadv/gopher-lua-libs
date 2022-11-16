local inspect = require 'inspect'
local require = require 'require'
local assert = require 'assert'

local tests = {
    {
        name = "number",
        value = 1,
    },
    {
        name = "string",
        value = "foobar",
    },
    {
        name = "object",
        value = {},
    },
}

function TestAssertTrue(t)
    for _, tt in pairs(tests) do
        t:Run(tt.name, function(t)
            assert:True(t, tt.value)
        end)
    end
end

function TestAssertFalse(t)
    for _, tt in pairs(tests) do
        t:Run(tt.name, function(t)
            assert:False(t, not tt.value)
        end)
    end
end

function TestAssertEqual(t)
    for _, tt in pairs(tests) do
        t:Run(tt.name, function(t)
            assert:Equal(t, tt.value, tt.value)
        end)
    end
end

function TestAssertNotEqual(t)
    for _, tt in pairs(tests) do
        t:Run(tt.name, function(t)
            assert:NotEqual(t, "DON'T MATCH", tt.value)
        end)
    end
end

function TestObjectsInspectEqual(t)
    local t1 = {
        foo = {
            bar = { "baz", "buz", "biz" }
        }
    }

    local t2 = {
        foo = {
            bar = { "baz", "buz", "biz" }
        }
    }
    require:Equal(t, inspect(t1), inspect(t2))
end