local http = require 'http'
local plugin = require 'plugin'
local time = require 'time'
local inspect = require 'inspect'

function TestMTLSServerWithClient(t)
    local addr_ch = channel.make(1)
    local server_body = [[
        local addr_ch = unpack(arg)
        local http = require 'http'
        local server, err = http.server {
            server_public_cert_pem_file = "test/data/test.cert.pem",
            server_private_key_pem_file = "test/data/test.key.pem",
        }
        assert(not err, tostring(err))
        addr_ch:send(server:addr())
        while true do
            local request, response = server:accept()
            response:code(200)
            response:write("OK\n")
            response:done()
        end
    ]]
    local server_plugin = plugin.do_string(server_body, addr_ch)
    server_plugin:run()
    time.sleep(1)
    local server_plugin_error = server_plugin:error()
    assert(not server_plugin_error, tostring(server_plugin_error))
    local ok, addr = addr_ch:receive(1)
    assert(ok, "addr not ok")
    local tURL = string.format("https://%s/", addr)

    t:Run('no-client-cert fails', function(t)
        local client = http.client()
        local req, err = http.request("GET", tURL)
        assert(not err, tostring(err))
        local resp, err = client:do_request(req)
        assert(err, tostring(err))
    end)

    t:Run('client-cert passes', function(t)
        local client = http.client {
            root_cas_pem_file = 'test/data/test.cert.pem',
            client_public_cert_pem_file = 'test/data/test.cert.pem',
            client_private_key_pem_file = 'test/data/test.key.pem',
        }
        local req, err = http.request("GET", tURL)
        assert(not err, tostring(err))
        local resp, err = client:do_request(req)
        assert(not err, tostring(err))
        assert(resp.code == 200, tostring(resp.code))
    end)

    server_plugin:stop()
end
