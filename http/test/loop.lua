local inspect = require("inspect")
print(inspect(request))
print(inspect(response))

print("host:", request.host)
print("method:", request.method)
print("referer:", request.referer)
print("proto:", request.proto)
print("request_uri:", request.request_uri)
print("remote_addr:", request.remote_addr)
for k, v in pairs(request.headers) do
  print("header: ", k, v)
end
-- get body
local body, err = request.body()
if err then error(err) end
print("body:", body)

response:code(200) -- write header
response:write(request.request_uri)
response:done()
