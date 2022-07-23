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

function TestReader(t)
    local s = [[{"foo":"bar","baz":"buz"}]]
    local reader = strings.new_reader(s)
    local got = reader:read("*a")
    assert(got == s, string.format("'%s' ~= '%s'", got, s))
end

function TestBuilder(t)
    local builder = strings.new_builder()
    builder:write("foo", "bar", 123)
    local got = builder:string()
    assert(got == "foobar123", string.format("'%s' ~= '%s'", got, "foobar123"))
end

function TestFields(t)
    local fields = strings.fields("a b c d")
    assert(#fields == 4, string.format("%d ~= 4", #fields))
    assert(fields[1] == "a", string.format("%s ~= 'a'", fields[1]))
    assert(fields[2] == "b", string.format("%s ~= 'b'", fields[2]))
    assert(fields[3] == "c", string.format("%s ~= 'c'", fields[3]))
    assert(fields[4] == "d", string.format("%s ~= 'd'", fields[3]))
end