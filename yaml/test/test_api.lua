local yaml = require("yaml")
local text = [[
a:
  b: 1
]]
local result, err = yaml.decode(text)
if err then error(err) end
if not(result["a"]["b"] == 1) then error("decode error") end

print("done: yaml.decode()")
