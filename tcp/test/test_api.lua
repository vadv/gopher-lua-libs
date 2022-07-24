local tcp = require("tcp")

function Test_tcp(t)

    local conn, err = tcp.open(":12345")
    assert(not err, err)
    t:Log("done: tcp:open()")

    local function assert_equal(expected, got)
        assert(got == expected, string.format([[expected "%s": got "%s"]], expected, got))
    end

    t:Run("write ping read pong", function(t)
        err = conn:write("ping")
        assert(not err, err)
        t:Log("done: tcp_client_ud:write()")

        local result, err = conn:read("*l")
        assert(not err, err)
        assert_equal("pong", result)
        t:Log("done: tcp_client_ud:read_line()")
    end)

    t:Run("read timeout fields", function(t)
        assert_equal(5, conn.dialTimeout)
        assert_equal(1, conn.writeTimeout)
        assert_equal(1, conn.readTimeout)
        assert_equal(1, conn.closeTimeout)
    end)

    t:Run('write/read timeout fields', function(t)
        -- Check setting fields
        conn.closeTimeout = 2
        assert_equal(2, conn.closeTimeout)
        conn.closeTimeout = 0.5
        assert_equal(0.5, conn.closeTimeout)
    end)
end
