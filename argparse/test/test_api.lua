local inspect = require "inspect"
local argparse = require "argparse"

local function assert_equal (expected, got)
    assert(got == expected, string.format([[expected "%s"; got "%s"]], expected, got))
end

local function assert_inspect_equal (expected, got)
    assert_equal(inspect(expected), inspect(got))
end

function Test_argparse(t)
    local parser = argparse("script", "An example.")
    parser:argument("input", "Input file.")
    parser:option("-o --output", "Output file.", "a.out")
    parser:option("-I --include", "Include locations."):count("*")

    local test_input = 'test_input'
    local test_output = 'test_input'
    local test_include = { "foo", "bar", "baz" }
    local test_arg = {
        "-o", test_output,
        test_input
    }
    for _, include in ipairs(test_include) do
        table.insert(test_arg, '-I')
        table.insert(test_arg, include)
    end

    local args = parser:parse(test_arg)
    print(inspect(args))
    assert_equal(test_input, args.input)
    assert_equal(test_output, args.output)
    assert_inspect_equal(test_include, args.include)
end

function Test_argparse_default(t)
    local parser = argparse("script", "An example.")
    parser:option("-f --foo", "Foo parameter", "default")
    local args = parser:parse({})
    assert(args.foo == 'default')
end

function Test_argparse_help(t)
    local parser = argparse("script", "An example.")
    parser:option("-f --foo", "Foo parameter", "default")
    local help = parser:get_help()
    assert(#help > 0, "help is empty")
end