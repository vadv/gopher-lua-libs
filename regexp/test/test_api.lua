local regexp = require("regexp")

local reg, err = regexp.compile("(gopher){2}")
if err then error(err) end
if reg:match("gopher") then error("must not be matched") end
if not reg:match("gophergopher") then error("must be matched") end
print("done: regexp_ud:match()")

local reg, err = regexp.compile("string: (.*)$")
if err then error(err) end
local result = reg:find_all_string_submatch("my string: 'hello world'")
if not(result[1][2] == "'hello world'") then error("not found: "..tostring(result[1][2])) end
print("done: regexp_ud:find_all_string_submatch():1")

local reg, err = regexp.compile("string: '(.*)\\s+(.*)'$")
if err then error(err) end
local result = reg:find_all_string_submatch("my string: 'hello world'")
if not(result[1][2] == "hello") then error("not found: "..tostring(result[1][2])) end
if not(result[1][3] == "world") then error("not found: "..tostring(result[1][3])) end
print("done: regexp_ud:find_all_string_submatch():2")

local found, err = regexp.match("(gopher){2}", "gophergopher")
if err then error(err) end
if not found then error("must be matched") end
print("done: regexp.match()")

local result, err = regexp.find_all_string_submatch("string: '(.*)\\s+(.*)'$", "my string: 'hello world'")
if err then error(err) end
if not(result[1][2] == "hello") then error("not found: "..tostring(result[1][2])) end
if not(result[1][3] == "world") then error("not found: "..tostring(result[1][3])) end
print("done: regexp.find_all_string_submatch()")
