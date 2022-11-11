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
    assert:False(true, "oh noes")
end