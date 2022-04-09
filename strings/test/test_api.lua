local strings = require("strings")

local str = "hello world"

function TestSplit(t)
    local str_parts = strings.split(str, " ")
    assert(type(str_parts) == 'table')
    assert(#str_parts == 2, string.format("%d ~= 2", #str_parts))
    assert(str_parts[1] == "hello", string.format("%s ~= hello", str_parts[1]))
    assert(str_parts[2] == "world", string.format("%s ~= world", str_parts[2]))
end

function TestHasPrefix(t)
    assert(strings.has_prefix(str, "hello"), [[not strings.has_prefix("hello")]])
end

function TestHasSuffix(t)
    assert(strings.has_suffix(str, "world"), [[not strings.has_suffix("world")]])
end

function TestTrim(t)
    assert(strings.trim(str, "world") == "hello ", "strings.trim()")
    assert(strings.trim(str, "hello ") == "world", "strings.trim()")
    assert(strings.trim_prefix(str, "hello ") == "world", "strings.trim()")
    assert(strings.trim_suffix(str, "hello ") == "hello world", "strings.trim()")
end

function TestContains(t)
    assert(strings.contains(str, "hello ") == true, "strings.contains()")
end