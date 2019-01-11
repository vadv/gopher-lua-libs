# http [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/http?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/http)

## Usage

### Client

```lua
local http = require("http")
local client = http.client()

-- GET
local request = http.request("GET", "http://hostname.com")
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code") end
if not(result.body == "xxx.xxx.xxx.xxx") then error("body") end

-- auth basic
local request = http.request("GET", "http://hostname.com")
request:set_basic_auth("admin", "123456")

-- headers
local client = http.client()
local request = http.request("POST", "http://hostname.com/api.json", "{}")
request:header_set("Content-Type", "application/json")

-- with proxy
local client = http.client({proxy="http(s)://login:password@hostname.com"})
local request = http.request("POST", "http://hostname.com/api.json", "{}")

-- ignore ssl
local client = http.client({insecure_ssl=true})
local request = http.request("POST", "http://hostname.com/api.json", "{}")

-- set headers for all request
local client = http.client({ headers={key="value"} })

-- set basic auth for all request
local client = http.client({basic_auth_user="admin", basic_auth_password="123456"})
```

### Server

```lua
local server, err = http.server("127.0.0.1:1113")
if err then error(err) end

while true do
  local req, resp = server:accept() -- lock and wait request

  -- print request
  print("host:", req.host)
  print("method:", req.method)
  print("referer:", req.referer)
  print("proto:", req.proto)
  print("path:", req.path)
  print("raw_path:", req.raw_path)
  print("raw_query:", req.raw_query)
  print("request_uri:", req.request_uri)
  print("remote_addr:", req.remote_addr)
  print("user_agent: "..req.user_agent)
  for k, v in pairs(req.headers) do
    print("header: ", k, v)
  end
  for k, v in pairs(req.query) do
    print("query params: ", k, "=" ,v)
  end
  -- write response
  resp:code(200) -- write header
  resp:header("content-type", "application/json")
  resp:write(req.request_uri) -- write data
  -- resp:redirect("http://google.com")
  resp:done() -- end response

end
```
