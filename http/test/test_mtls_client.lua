local http = require 'http_client'

function TestMTLS(t)
    assert(tURL, 'tURL global is not set')

    t:Run('no-client-cert fails', function(t)
        local client = http.client{
            insecure_ssl = true,
        }
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
end
