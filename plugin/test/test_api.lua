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

local curl_plugin = plugin.new(plugin_body_1)

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
