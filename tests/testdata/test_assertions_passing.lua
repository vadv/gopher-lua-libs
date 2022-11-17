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

function TestNil(t)
    assert:Nil(t, nil)
end

function TestNotNil(t)
    assert:NotNil(t, 123)
    assert:NotNil(t, "")
end


function TestGreater(t)
    assert:Greater(t, 4, 1)
end

function TestGreaterf(t)
    assert:Greaterf(t, 4, 1, "foo %s", "bar")
end

function TestGreaterOrEqual(t)
    assert:GreaterOrEqual(t, 4, 4)
end

function TestGreaterOrEqualf(t)
    assert:GreaterOrEqualf(t, 4, 4, "foo %s", "bar")
end

function TestLess(t)
    assert:Less(t, 2, 3)
end

function TestLessf(t)
    assert:Less(t, 2, 3, "abc %s", "def")
end

function TestLessOrEqual(t)
    assert:LessOrEqual(t, 2, 2)
end

function TestLessOrEqualf(t)
    assert:LessOrEqualf(t, 2, 2, "abc %s", "def")
end
