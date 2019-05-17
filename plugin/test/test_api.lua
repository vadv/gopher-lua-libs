local plugin = require("plugin")
local time = require("time")
local ioutil = require("ioutil")

local plugin_body_1 = [[

local http = require("http")
local time = require("time")

local i = 1
local client = http.client()
local request = http.request("GET", "http://google.com")
os.remove("./test/file.txt")

while true do
    local result, err = client:do_request(request)
    if err then error(err) end
    time.sleep(0.1)
    local file = io.open("./test/file.txt", "w")
    file:write(tostring(i), "\n")
    file:close()
    i = i + 1
end

]]

local curl_plugin = plugin.do_string(plugin_body_1)

curl_plugin:run()
time.sleep(2)
curl_plugin:stop()
time.sleep(1)


local running = curl_plugin:is_running()
if running then error("already running") end

local data, err = ioutil.read_file("./test/file.txt")
if err then error(err) end

print(data)
local i = tonumber(data)
if i < 1 then error("i < 1") end


-- test with payload
local plugin_body_2 = [[
    local ioutil = require("ioutil")
    local err = ioutil.write_file("./test/payload.txt", payload)
    if err then error(err) end
]]
local test_payload = "OK"
local plugin_with_payload = plugin.do_string_with_payload(plugin_body_2, test_payload)
plugin_with_payload:run()

time.sleep(1)
local i = 0
while plugin_with_payload:is_running() do
    time.sleep(1)
    i = i + 1
    if i > 3 then error("timeout") end
end

local data, err = ioutil.read_file("./test/payload.txt")
if err then error(err) end
if not(data == test_payload) then error("data <> test_payload") end
