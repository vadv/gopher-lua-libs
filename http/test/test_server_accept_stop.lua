local http = require("http")
local plugin = require("plugin")
local time = require("time")

local plugin_body = [[
  local http = require("http_server")
  local server, err = http.server("127.0.0.1:3113")
  if err then error(err) end
  while true do
    local request, response = server:accept()

    -- raise internal error
    if request.request_uri == "/bug" then
      error("get /bug")
    end

    response:code(200) -- write header
    response:write(request.request_uri)
    response:done()
  end
]]

local p = plugin.do_string(plugin_body)
p:run()
time.sleep(1)

local client = http.client({timeout=1})
local req, err = http.request("GET", "http://127.0.0.1:3113")
if err then error(err) end

-- no error request
local _, err = client:do_request(req)
if err then error(err) end

-- stop plugin
p:stop()
time.sleep(1)
if p:error() then error(err) end

-- must error request
local _, err = client:do_request(req)
if (err == nil) then error("must be error") end

p:run()
time.sleep(1)
if p:error() then error( p:error() ) end
-- no error request
local _, err = client:do_request(req)
if err then error(err) end

-- raise internal error
local req_bug, err = http.request("GET", "http://127.0.0.1:3113/bug")
if err then error(err) end
local _, err = client:do_request(req_bug)
if not(err) then error("must be internal error") end
if not(p:error()) then error("must be internal error") end

-- test successful start
p:stop()
time.sleep(1)
p:run()
time.sleep(1)
if p:error() then error( p:error() ) end
-- no error request
local _, err = client:do_request(req)
if err then error(err) end