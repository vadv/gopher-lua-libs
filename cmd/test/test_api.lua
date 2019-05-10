local cmd = require("cmd")
local runtime = require("runtime")

local command = "sleep 1"
if runtime.goos() == "windows" then command = "timeout 1" end

local result, err = cmd.exec(command)
if err then error(err) end
print(result.stdout)

-- Test timeout
local cmd = require("cmd")
local runtime = require("runtime")

local command = "sleep 5"
if runtime.goos() == "windows" then command = "timeout 1" end

local result, err = cmd.exec(command, 1)
if err == nil then error("timeout expected but did not occur") end
if err ~= "execute timeout" then error("expected 'execute timeout' but instead got '" .. err .. "'") end