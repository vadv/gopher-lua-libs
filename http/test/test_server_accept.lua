local http = require("http_server")

local server, err = http.server("127.0.0.1:1113")
if err then error(err) end

local running, count = true, 0
while running do
  local request, response = server:accept()
  print("host:", request.host)
  print("method:", request.method)
  print("referer:", request.referer)
  print("proto:", request.proto)
  print("request_uri:", request.request_uri)
  print("remote_addr:", request.remote_addr)
  for k, v in pairs(request.headers) do
    print("header: ", k, v)
  end
  response:code(200) -- write header
  response:write(request.request_uri)
  response:done()
  count = count + 1
  running = (count < 10)
end

if count < 10 then
  error("count: "..tostring(count))
end
