local shellescape = require("shellescape")
local inspect = require('inspect')

local test = {}

-- Test string without quotes is string
function test:quote_foo_is_foo()
    input = "foo"
    expected = "foo"
    escaped = shellescape.quote(input)
    assert(escaped == expected, string.format("%s != %s", escaped, expected))
end

-- Test string with space is quoted
function test:quote_text_with_spaces_is_quoted()
    input = "foo bar"
    expected = "'foo bar'"
    escaped = shellescape.quote(input)
    assert(escaped == expected, string.format("%s != %s", escaped, expected))
end

-- Test string with apostrophe is quoted
function test:quote_text_with_apostrophe_is_quoted()
    input = "don't sweat the technique"
    expected = [['don'"'"'t sweat the technique']]
    escaped = shellescape.quote(input)
    assert(escaped == expected, string.format("%s != %s", escaped, expected))
end

-- Test string with double quotes is quoted
function test:quote_text_with_double_quotes_is_quoted()
    input = [[She said, "that's what she said."]]
    expected = [['She said, "that'"'"'s what she said."']]
    escaped = shellescape.quote(input)
    assert(escaped == expected, string.format("%s != %s", escaped, expected))
end

-- Test QuoteCommand: mixed args are all quoted
function test:quote_command__array_with_mixture_is_escaped()
    input = {"foo", "bar baz", "don't do it"}
    expected = [[foo 'bar baz' 'don'"'"'t do it']]
    escaped = shellescape.quote_command(input)
    assert(escaped == expected, string.format("%s != %s", escaped, expected))
end

-- Test StripUnsafe: safe chars are left alone
function test:strip_unsafe__safe_are_preserved()
    input = [[don't say "foo."]]
    expected = [[don't say "foo."]]
    escaped = shellescape.strip_unsafe(input)
    assert(escaped == expected, string.format("%s != %s", escaped, expected))
end

-- Test StripUnsafe: safe chars are left alone
function test:strip_unsafe__remove_newline()
    input = "test\ntwolines"
    expected = "testtwolines"
    escaped = shellescape.strip_unsafe(input)
    assert(escaped == expected, string.format("%s != %s", escaped, expected))
end

return test
