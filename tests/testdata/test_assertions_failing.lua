local inspect = require 'inspect'
local require = require 'require'
local assert = require 'assert'

function TestAssertGivesTwoResults(t)
    assert:Equal(t, "one", "two", "foobar")
    assert:Equalf(t, "three", "four", "This is a test %s: %d", "foobar", 123)
end

function TestRequireThenAssertGivesOneResult(t)
    require:Equal(t, "abc", "def")
    assert:Equal(t, true, false, "shouldn't reach here")
end

function TestTrue(t)
    assert:True(t, false, 'zoiks')
    assert:Truef(t, false, 'crikey %d', 123)
end

function TestTablesEqual(t)
    local t1 = {
        foo = "bar"
    }
    local t2 = {
        foo = "notbar"
    }
    assert:Equal(t, inspect(t1), inspect(t2, { process = remove_all_metatables }))
end

function TestTablesNotEqual(t)
    local t1 = {
        foo = "bar"
    }
    local t2 = {
        foo = "bar"
    }
    assert:NotEqual(t, inspect(t1), inspect(t2, { process = remove_all_metatables }))
end

function TestFalse(t)
    assert:False(t, true, "oh noes")
end

function TestStringsEqual(t)
    assert:Equal(t, "expected text", "actual text")
end

function TestNoError(t)
    assert:NoError(t, "foo bar")
end

function TestNil(t)
    assert:Nil(t, "foo bar baz")
end

function TestNotNil(t)
    assert:NotNil(t, nil)
end

function TestGreater(t)
    assert:Greater(t, 4, 5)
end

function TestGreaterf(t)
    assert:Greaterf(t, 4, 5, "foo %s", "bar")
end

function TestGreaterOrEqual(t)
    assert:GreaterOrEqual(t, 4, 5)
end

function TestGreaterOrEqualf(t)
    assert:GreaterOrEqualf(t, 4, 5, "foo %s", "bar")
end

function TestLess(t)
    assert:Less(t, 2, 1)
end

function TestLessf(t)
    assert:Lessf(t, 2, 1, "abc %s", "def")
end

function TestLessOrEqual(t)
    assert:LessOrEqual(t, 2, 1)
end

function TestLessOrEqualf(t)
    assert:LessOrEqualf(t, 2, 1, "abc %s", "def")
end
