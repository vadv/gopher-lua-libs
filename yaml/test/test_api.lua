local yaml = require("yaml")
local test = {}

-- test decode
function Test_decode(t)
    local text = [[
a:
  b: 1
]]
    local result, err = yaml.decode(text)
    assert(not err, tostring(err))
    assert(result["a"]["b"] == 1, tostring(result["a"]["b"]))
    print("done: yaml.decode(t)n")
end

-- test decode with no args throws exception
function Test_decode_no_args(t)
    local ok, errMsg = pcall(yaml.decode)
    assert(not ok)
    assert(errMsg)
    assert(errMsg:find("(string expected, got nil)"), tostring(errMsg))
end

-- test encode of decoded(text) == text
function Test_decoded_text_equals_text(t)
    local text = [[
a:
  b: 1
]]
    local result = { a = { b = 1 } }
    local encodedResult, err = yaml.encode(result)
    assert(not err, tostring(err))
    assert(encodedResult == text, tostring(encodedResult)
    )
end

-- test encode(slice) works
function Test_encode_slice_works(t)
    local encodedSlice = yaml.encode({ "foo", "bar", "baz" })
    assert(encodedSlice == [[
- foo
- bar
- baz
]], tostring(encodedSlice))
end

-- test encode(sparse slice) works
function Test_encode_sparse_slice_returns_map(t)
    local slice = { [0] = "foo", [1] = "bar", [2] = "baz" }
    local encodedSlice = yaml.encode(slice)
    assert(encodedSlice == [[
0: foo
1: bar
2: baz
]], tostring(encodedSlice))
end

-- test encode(map) works
function Test_encode_map_returns_map(t)
    local map = { foo = "bar", bar = { 1, 2, 3.45 } }
    local encodedMap = yaml.encode(map)
    assert(encodedMap == [[
bar:
- 1
- 2
- 3.45
foo: bar
]], tostring(encodedMap))
end

-- test encode(function) fails
function Test_encode_function_fails(t)
    local _, errMsg = yaml.encode(function()
        return ""
    end)
    assert(errMsg)
    assert(errMsg:find("cannot encode values with function in them"), errMsg)

    -- test encode with no args throws exception
    local ok, errMsg = pcall(yaml.encode)
    assert(not ok)
    assert(errMsg:find("(value expected)"), tostring(errMsg))
end

-- test cycles
function Test_cycles_return_error(t)
    local t1 = {}
    local t2 = { t1 = t1 }
    t1[t2] = t2
    local _, errMsg = yaml.encode(t1)
    assert(errMsg)
    assert(errMsg:find("nested table"), tostring(errMsg))
end

return test
