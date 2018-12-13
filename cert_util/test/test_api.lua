local cert_util = require("cert_util")

local tx, err = cert_util.not_after("127.0.0.1:1443", "127.0.0.1:1443")
if err then error(err) end

if not(tx == 1576161031) then error("cert") end
