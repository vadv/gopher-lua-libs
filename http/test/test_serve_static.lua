local http = require("http")
local plugin = require("plugin")
local time = require("time")

local plugin_body = [[
local http = require("http")
local err = http.serve_static("./test/data", "127.0.0.1:2115")
if err then error(err) end
]]

local p = plugin.do_string(plugin_body)
p:run()
time.sleep(1)
if not p:is_running() then
    error(p:error())
end

local client = http.client()
local req, err = http.request("GET", "http://127.0.0.1:2115")
if err then error(err) end
local resp, err = client:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
if not(resp.body == "OK") then error("resp body, get:"..resp.body) end