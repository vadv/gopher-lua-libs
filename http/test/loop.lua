local inspect = require("inspect")
print(inspect(request))
print(inspect(response))

response:code(200) -- write header
response:write(request.request_uri)
response:done()
