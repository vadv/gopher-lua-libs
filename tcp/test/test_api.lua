local strings = require("strings")
local tcp = require("tcp")

function Test_tcp(t)

    local conn, err = tcp.open(":12345")
    assert(not err, err)
    t:Log("done: tcp:open()")

    err = conn:write("ping")
    assert(not err, err)
    t:Log("done: tcp_client_ud:write()")

    local result, err = conn:read()
    assert(not err, err)
    assert(strings.trim_space(result) == "pong", string.format([[expected "%s"; got "%s"]], "pong", result))
    t:Log("done: tcp_client_ud:read_line()")
end