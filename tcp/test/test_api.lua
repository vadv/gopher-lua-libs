local tcp = require("tcp")

local conn, err = tcp.open(":12345")
if err then error(err) end
print("done: tcp:open()")

err = conn:write("ping")
if err then error(err) end
print("done: tcp_client_ud:write()")

local result, err = conn:read()
if err then error(err) end
if (result ~= "pong\n") then error("must be pong message") end
print("done: tcp_client_ud:read_line()")
