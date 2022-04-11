local json = require("json")
local inspect = require("inspect")
local io = require("io")

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
