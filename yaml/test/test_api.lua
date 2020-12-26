-- test decode
local yaml = require("yaml")
local text = [[
a:
  b: 1
]]
local result, err = yaml.decode(text)
assert(not err, tostring(err))
assert(result["a"]["b"] == 1, tostring(result["a"]["b"]))
print("done: yaml.decode()")

-- test decode with no args throws exception
local ok, errMsg = pcall(yaml.decode)
assert(not ok)
assert(errMsg)
assert(errMsg:find("(string expected, got nil)"), tostring(errMsg))

-- test encode of decoded(text) == text
local encodedResult, err = yaml.encode(result)
assert(not err, tostring(err))
assert(encodedResult == text, tostring(encodedResult))

-- test encode(slice) works
local encodedSlice = yaml.encode({ "foo", "bar", "baz" })
assert(encodedSlice == [[
- foo
- bar
- baz
]], tostring(encodedSlice))

-- test encode(sparse slice) works
local slice = { [0] = "foo", [1] = "bar", [2] = "baz" }
local encodedSlice = yaml.encode(slice)
assert(encodedSlice == [[
0: foo
1: bar
2: baz
]], tostring(encodedSlice))

-- test encode(map) works
local map = { foo = "bar", bar = { 1, 2, 3.45 } }
local encodedMap = yaml.encode(map)
assert(encodedMap == [[
bar:
- 1
- 2
- 3.45
foo: bar
]], tostring(encodedMap))

-- test encode(function) fails
local _, errMsg = yaml.encode(function() return "" end)
assert(errMsg)
assert(errMsg:find("cannot encode values with function in them"), errMsg)

-- test encode with no args throws exception
local ok, errMsg = pcall(yaml.encode)
assert(not ok)
assert(errMsg:find("(value expected)"), tostring(errMsg))
