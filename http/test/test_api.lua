local http = require("http")

local client_1 = http.client()
local client_2 = http.client({timeout=1})
local client_3 = http.client({insecure_skip_verify=true})
local client_4 = http.client({insecure_skip_verify=false, timeout=1})
local client_5 = http.client({basic_auth_user="admin", basic_auth_password="123456"})
local client_6 = http.client({headers={simple_header="check_header"}})
local client_7 = http.client({user_agent="check_ua"})

local req, err = http.request("GET", "http://127.0.0.1:1111/get")
if err then error(err) end
local resp, err = client_1:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
if not(resp.body == "OK") then error("resp body") end
print("done: http.client:get")

local req, err = http.request("GET", "http://127.0.0.1:1111/getBasicAuth")
if err then error(err) end
req:set_basic_auth("admin", "123456")
local resp, err = client_1:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
if not(resp.body == "OK") then error("resp body") end
print("done: http.client:getBasicAuth")

local req, err = http.request("GET", "http://127.0.0.1:1111/timeout")
if err then error(err) end
local _, err = client_2:do_request(req)
if err == nil then error("must be error") end
print("done: http.client:timeout")

local req, err = http.request("GET", "https://127.0.0.1:1112/get")
if err then error(err) end
local _, err = client_3:do_request(req)
if err == nil then error("must be error") end
print("done: http.client:ssl+error")

local req, err = http.request("GET", "http://127.0.0.1:1111/getBasicAuth")
if err then error(err) end
req:set_basic_auth("admin", "123456")
local resp, err = client_4:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
if not(resp.body == "OK") then error("resp body") end
print("done: http.client:ssl+getBasicAuth")

local req, err = http.request("GET", "https://127.0.0.1:1112/timeout")
if err then error(err) end
local _, err = client_4:do_request(req)
if err == nil then error("must be error") end
print("done: http.client:ssl+timeout")

local test_unescape = "<> dasdsadas"
local test_escape = "%3C%3E+dasdsadas"
if not (http.query_escape(test_unescape) == test_escape) then error("escape error") end
print("done: http.escape")

local req, err = http.request("GET", "http://127.0.0.1:1111/getBasicAuth")
if err then error(err) end
local resp, err = client_5:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
if not(resp.body == "OK") then error("resp body") end
print("done: http.client:getBasicAuth via client")

local req, err = http.request("GET", "http://127.0.0.1:1111/checkHeader")
if err then error(err) end
local resp, err = client_6:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
if not(resp.body == "OK") then error("resp body") end
print("done: http.client:header via client")

local req, err = http.request("GET", "http://127.0.0.1:1111/checkUserAgent")
if err then error(err) end
local resp, err = client_7:do_request(req)
if err then error(err) end
if not(resp.code == 200) then error("resp code") end
if not(resp.body == "OK") then error("resp body") end
print("done: http.client:user agent via client")

local server, err = http.server("127.0.0.1:1113")
if err then error(err) end

local running, count = true, 0
while running do
  local req, response = server:accept()
  print("host:", req.host)
  print("method:", req.method)
  print("referer:", req.referer)
  print("proto:", req.proto)
  print("request_uri:", req.request_uri)
  print("remote_addr:", req.remote_addr)
  for k, v in pairs(req.headers) do
    print("header: ", k, v)
  end
  response:code(200) -- write header
  response:write(req.request_uri)
  response:done()
  count = count + 1
  running = (count < 10)
end
