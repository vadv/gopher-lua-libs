local http = require("http_server")

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
