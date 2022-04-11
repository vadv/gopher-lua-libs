local json = require("json")
local inspect = require("inspect")
local io = require("io")
local strings = require("strings")

function TestJson(t)
    local jsonStringWithNull = [[{"a":{"b":1, "c":null}}]]
    local jsonString = [[{"a":{"b":1}}]]

    local result, err = json.decode(jsonStringWithNull)
    t:Run("decode", function(t)
        if err then error(err) end
        if result["a"]["b"] ~= 1 then error("must be decode") end
        if result["a"]["c"] ~= nil then error("c is not nil") end
        print("done: json.decode()")
    end)

    local result, err = json.encode(result)
    t:Run("encode omits null values", function(t)
        if err then error(err) end
        if result ~= jsonString then error("must be encode "..inspect(result)) end
        print("done: json.encode()")
    end)
end

function TestEncoder(t)
    temp_file = '/tmp/tst.json'
    os.remove(temp_file)
    writer, err = io.open(temp_file, 'w')
    assert(not err, err)
    encoder = json.new_encoder(writer)
    err = encoder:encode({foo="bar", bar="baz"})
    assert(not err, err)
    writer:close()

    reader = io.open(temp_file, 'r')
    contents = reader:read('*a')
    assert(contents, "contents should not be empty")
    contents = json.decode(contents)
    assert(contents['foo'] == 'bar', string.format("%s ~= bar", contents['foo']))
    assert(contents['bar'] == 'baz', string.format("%s ~= baz", contents['bar']))
end

function TestEncoderWithStringsBuffer(t)
    builder = strings.new_builder()
    encoder = json.new_encoder(builder)
    err = encoder:encode({abc="def", num=123, arr={1,2,3}})
    s = strings.trim_suffix(builder:string(), "\n")
    expected = [[{"abc":"def","arr":[1,2,3],"num":123}]]
    assert(s == expected, string.format([['%s' ~= '%s']], expected, s))
end

function TestEncoderWithPrettyPrinting(t)
    builder = strings.new_builder()
    encoder = json.new_encoder(builder)
    encoder:set_indent('', "  ")
    err = encoder:encode({abc="def", num=123, arr={1,2,3}})
    s = strings.trim_suffix(builder:string(), "\n")
    expected = [[{
  "abc": "def",
  "arr": [
    1,
    2,
    3
  ],
  "num": 123
}]]
    assert(s == expected, string.format([['%s' ~= '%s']], expected, s))
end

function TestDecoder(t)
    temp_file = '/tmp/tst.json'
    os.remove(temp_file)
    writer, err = io.open(temp_file, 'w')
    assert(not err, err)
    writer:write([[{"abc": "def", "num": 123}]])
    writer:close()

    reader = io.open(temp_file, 'r')
    decoder = json.new_decoder(reader)
    result, err = decoder:decode()
    assert(not err, err)
    assert(result['abc'] == 'def', string.format("%s ~= def", result['abc']))
    assert(result['num'] == 123, string.format("%d ~= 123", result['num']))
end

function TestDecoderWithStringsReader(t)
    s = [[{"abc": "def", "num": 123}]]
    reader = strings.new_reader(s)
    decoder = json.new_decoder(reader)
    result, err = decoder:decode()
    assert(not err, err)
    assert(result['abc'] == 'def', string.format("%s ~= def", result['abc']))
    assert(result['num'] == 123, string.format("%d ~= 123", result['num']))
end

function TestDecoder_reading_twice(t)
    input = [[
{"abc": "def"}
{"num": 123}
]]
    reader = strings.new_reader(input)
    decoder = json.new_decoder(reader)
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
    encoder = json.new_encoder(writer)
    err = encoder:encode({abc="def"})
    assert(not err, err)
    encoder:encode({num=123})
    assert(not err, err)
    s = writer:string()
    expected = [[{"abc":"def"}
{"num":123}
]]
    assert(s == expected, string.format([['%s' ~= '%s']], s, expected))
end
