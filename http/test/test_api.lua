local http = require("http")
local client_1 = http.client()
local client_2 = http.client({timeout=1})
local client_3 = http.client({insecure_skip_verify=true})
local client_4 = http.client({insecure_skip_verify=false, timeout=1})

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

local req, err = http.request("GET", "http://127.0.0.1:1112/get")
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

local req, err = http.request("GET", "http://127.0.0.1:1112/timeout")
if err then error(err) end
local _, err = client_4:do_request(req)
if err == nil then error("must be error") end
print("done: http.client:ssl+timeout")

local test_unescape = "<> dasdsadas"
local test_escape = "%3C%3E+dasdsadas"
if not (http.query_escape(test_unescape) == test_escape) then error("escape error") end
print("done: http.escape")
