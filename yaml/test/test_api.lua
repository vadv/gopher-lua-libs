local yaml = require("yaml")
local io = require("io")
local strings = require("strings")

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

function TestEncoder(t)
    temp_file = '/tmp/tst.json'
    os.remove(temp_file)
    writer, err = io.open(temp_file, 'w')
    assert(not err, err)
    encoder = yaml.new_encoder(writer)
    err = encoder:encode({foo="bar", bar="baz"})
    assert(not err, err)
    writer:close()

    reader = io.open(temp_file, 'r')
    contents = reader:read('*a')
    assert(contents, "contents should not be empty")
    contents = yaml.decode(contents)
    assert(contents['foo'] == 'bar', string.format("%s ~= bar", contents['foo']))
    assert(contents['bar'] == 'baz', string.format("%s ~= baz", contents['bar']))
end

function TestEncoderWithStringsBuffer(t)
    builder = strings.new_builder()
    encoder = yaml.new_encoder(builder)
    err = encoder:encode({abc="def", num=123, arr={1,2,3}})
    s = strings.trim(builder:string(), "\n")
    expected = strings.trim([[
abc: def
arr:
- 1
- 2
- 3
num: 123
]], " \n")
    assert(s == expected, string.format([['%s' ~= '%s']], expected, s))
end

function TestDecoder(t)
    temp_file = '/tmp/tst.json'
    os.remove(temp_file)
    writer, err = io.open(temp_file, 'w')
    assert(not err, err)
    writer:write([[
abc: def
num: 123
]])
    writer:close()

    reader = io.open(temp_file, 'r')
    decoder = yaml.new_decoder(reader)
    result, err = decoder:decode()
    assert(not err, err)
    assert(result['abc'] == 'def', string.format("%s ~= def", result['abc']))
    assert(result['num'] == 123, string.format("%d ~= 123", result['num']))
end

function TestDecoderWithStringsReader(t)
    s = [[
abc: def
num: 123
]]
    reader = strings.new_reader(s)
    decoder = yaml.new_decoder(reader)
    result, err = decoder:decode()
    assert(not err, err)
    assert(result['abc'] == 'def', string.format("%s ~= def", result['abc']))
    assert(result['num'] == 123, string.format("%d ~= 123", result['num']))
end

function TestDecoder_reading_twice(t)
    input = [[
abc: def
---
num: 123
]]
    reader = strings.new_reader(input)
    decoder = yaml.new_decoder(reader)
    first, err = decoder:decode()
    assert(not err, err)
    second, err = decoder:decode()
    assert(not err, err)

    s = first["abc"]
    expected = "def"
    assert(s == expected, string.format([['%s' ~= '%s']], s, expected))

    num = second["num"]
    expected = 123
    assert(num == expected, string.format([['%d' ~= '%d']], num, expected))
end

function TestEncoder_writing_twice(t)
    writer = strings.new_builder()
    encoder = yaml.new_encoder(writer)
    expected = "def"
    err = encoder:encode({abc="def"})
    assert(not err, err)
    err = encoder:encode({num=123})
    assert(not err, err)
    s = writer:string()
    expected = [[
abc: def
---
num: 123
]]
    assert(s == expected, string.format([['%s' ~= '%s']], s, expected))
end
