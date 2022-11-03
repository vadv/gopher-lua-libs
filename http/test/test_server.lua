local http = require 'http'
local plugin = require 'plugin'
local time = require 'time'
local inspect = require 'inspect'

function Test_do_handle_function(t)
    local addr_ch = channel.make(1)
    local server_body = [[
        local addr_ch = unpack(arg)
        local http = require 'http'
        local server, err = http.server {}
        assert(not err, tostring(err))
        addr_ch:send(server:addr())
        server:do_handle_function(function(response, request)
            print(string.format("response = %s", response))
            response:code(200)
            response:write("OK\n")
            response:done()
        end)
    ]]
    local server_plugin = plugin.do_string(server_body, addr_ch)
    server_plugin:run()
    time.sleep(1)
    local server_plugin_error = server_plugin:error()
    assert(not server_plugin_error, tostring(server_plugin_error))
    local ok, addr = addr_ch:receive(1)
    assert(ok, "addr not ok")
    local tURL = string.format("http://%s/", addr)

    local client = http.client()
    local req = http.request('GET', tURL)
    local resp, err = client:do_request(req)
    assert(not err, tostring(err))
    assert(resp.code == 200, tostring(resp.code))
end