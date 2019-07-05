local prometheus = require("prometheus")
local http = require("http_client")
local strings = require("strings")

local pp = prometheus.register(":18080")
pp:start()

local gauge = prometheus.gauge({
     namespace="node_scout",
     subsystem="nf_conntrack",
     name="insert_failed",
     help="insert_failed from nf_conntrack",
})
gauge:set(100)

local client = http.client({timeout=5})
local request = http.request("GET", "http://127.0.0.1:18080/metrics")

local result, err = client:do_request(request)
if err then erorr(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed 100")) then error("body:\n"..result.body) end

gauge:add(1)
local result, err = client:do_request(request)
if err then erorr(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed 101")) then error("body:\n"..result.body) end
