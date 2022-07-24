local pprof = require("pprof")
local http = require("http")
local time = require("time")

function Test_pprof(t)
    local client = http.client()
    local pp = pprof.register(":1234")

    pp:enable()
    time.sleep(1)

    local req, err = http.request("GET", "http://127.0.0.1:1234/debug/pprof/goroutine")
    assert(not err, err)
    local resp, err = client:do_request(req)
    assert(not err, err)
    assert(resp.code == 200, "resp code: " .. resp.code)

    pp:disable()
    time.sleep(5)

    local resp, err = client:do_request(req)
    assert(err, "must be error")
end
