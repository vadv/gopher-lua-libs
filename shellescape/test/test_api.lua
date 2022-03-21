local shellescape = require("shellescape")
local inspect = require('inspect')

-- Tests for the Quote function
function TestQuote(t)
    tests = {
        {
            name="simple string is itself",
            input="foo",
            expected="foo",
        },
        {
            name="string with spaces is quoted",
            input="foo bar",
            expected=[['foo bar']],
        },
        {
            name="string with apostrophe is quoted",
            input="don't sweat the technique",
            expected=[['don'"'"'t sweat the technique']],
        },
        {
            name="string with double quotes is quoted",
            input=[[She said, "that's what she said."]],
            expected=[['She said, "that'"'"'s what she said."']],
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            escaped = shellescape.quote(tt.input)
            assert(escaped == tt.expected, string.format("%s != %s", escaped, tt.expected))
        end)
    end

end

-- Tests for QuoteCommand
function TestQuoteCommand(t)
    tests = {
        {
            name="mixed args are all quoted",
            input={"foo", "bar baz", "don't do it"},
            expected=[[foo 'bar baz' 'don'"'"'t do it']],
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            escaped = shellescape.quote_command(tt.input)
            assert(escaped == tt.expected, string.format("%s != %s", escaped, tt.expected))
        end)
    end
end

-- Tests for StripUnsafe
function TestStripUnsafe(t)
    tests = {
        {
            name="safe chars are left alone",
            input=[[don't say "foo."]],
            expected=[[don't say "foo."]],
        },
        {
            name="unsafe chars like newline are removed",
            input="foo\nbar",
            expected="foobar",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            escaped = shellescape.strip_unsafe(tt.input)
            assert(escaped == tt.expected, string.format("%s != %s", escaped, tt.expected))
        end)
    end
end
