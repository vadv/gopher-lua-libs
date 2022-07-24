local tcp = require("tcp")

local conn, err = tcp.open(":12345")
if err then
    error(err)
end
print("done: tcp:open()")

err = conn:write("ping")
if err then
    error(err)
end
print("done: tcp_client_ud:write()")

local result, err = conn:read()
if err then
    error(err)
end
if (result ~= "pong\n") then
    error("must be pong message")
end
print("done: tcp_client_ud:read_line()")

-- Check fields
function assert_equal(expected, got)
    assert(got == expected, string.format("expected %s: got %s", expected, got))
end
assert_equal(5, conn.dialTimeout)
assert_equal(1, conn.writeTimeout)
assert_equal(1, conn.readTimeout)
assert_equal(1, conn.closeTimeout)
print("done: tcp_client_ud read timeout fields")

-- Check setting fields
conn.closeTimeout = 2
assert_equal(2, conn.closeTimeout)
conn.closeTimeout = 0.5
assert_equal(0.5, conn.closeTimeout)
print("done: tcp_client_ud write and check timeout fields")
