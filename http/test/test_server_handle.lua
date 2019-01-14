local http = require("http")
local plugin = require("plugin")
local time = require("time")

local plugin_body = [[
local http = require("http")
local server, err = http.server("127.0.0.1:2113")
if err then error(err) end
server:do_handle_file("./test/loop.lua")
]]

local p = plugin.do_string(plugin_body)
p:run()
time.sleep(5)
