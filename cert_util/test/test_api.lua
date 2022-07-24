local cert_util = require("cert_util")

function Test_cert_util(t)
    local tx, err = cert_util.not_after("127.0.0.1:1443", "127.0.0.1:1443")
    assert(not err, err)
    assert(tx == 1576161031, "cert: " .. tx)
end
