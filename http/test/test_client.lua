local http = require("http_client")
local http_util = require("http_util")

local client_1 = http.client()
local client_2 = http.client({timeout=1})
local client_3 = http.client({insecure_ssl=false})
local client_4 = http.client({insecure_ssl=true, timeout=1})
local client_5 = http.client({basic_auth_user="admin", basic_auth_password="123456"})
local client_6 = http.client({headers={simple_header="check_header"}})
local client_7 = http.client({user_agent="check_ua"})

local req, err = http.request("GET", "http://127.0.0.1:1111/get")
if err then error(err) end
local resp, err = client_1:do_request(req)
if err then error(err) end
if not(resp.headers['Content-Length'] == "2") then error("resp headers") end
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

local req, err = http.request("GET", "https://127.0.0.1:1112/getBasicAuth")
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
if not (http_util.query_escape(test_unescape) == test_escape) then error("escape error") end
print("done: http.escape")

local url, err = http_util.parse_url("http://user1:password11@127.0.0.1:1111/getBasicAuth?k1=v1&k2=x&k1=v2")
if err then error(err) end
if not(url.scheme == "http") then
  error("must be scheme http")
end
if not(url.host == "127.0.0.1:1111") then
  error("must be host 127.0.0.1:1111")
end
if not(url.user.username == "user1") then
  error("must be user1")
end
if not(url.user.password == "password11") then
  error("must be password11")
end
if not(url.query.k1[2] == "v2") then
  error("must be v2")
end
if not(url.path == "/getBasicAuth") then
  error("must be /getBasicAuth")
end

url.path = "/test"
if not(http_util.build_url(url) == "http://user1:password11@127.0.0.1:1111/test?k1=v1&k1=v2&k2=x") then
  error("get: "..http_util.build_url(url))
end

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
