local prometheus = require("prometheus")
local http = require("http_client")
local strings = require("strings")
local time = require("time")

local pp = prometheus.register(":18080")
pp:start()

time.sleep(1)

-- gauge
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
if err then error(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed 100")) then error("body:\n"..result.body) end

gauge:add(1)
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed 101")) then error("body:\n"..result.body) end

-- gauge vector
local gauge_vec, err = prometheus.gauge({
    namespace="node_scout",
    subsystem="nf_conntrack",
    name="insert_failed_vec",
    help="insert_failed from nf_conntrack",
    labels={"label_1", "label_2"},
})
if err then error(err) end

gauge_vec:set(100, {label_1="one", label_2="two"})
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, 'node_scout_nf_conntrack_insert_failed_vec{label_1="one",label_2="two"} 100')) then error("body:\n"..result.body) end

gauge_vec:add(1, {label_1="one", label_2="two"})
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, 'node_scout_nf_conntrack_insert_failed_vec{label_1="one",label_2="two"} 101')) then error("body:\n"..result.body) end

gauge_vec:inc({label_1="one", label_2="two"})
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, 'node_scout_nf_conntrack_insert_failed_vec{label_1="one",label_2="two"} 102')) then error("body:\n"..result.body) end

-- re-register gauge vector
local _, err = prometheus.gauge({
    namespace="node_scout",
    subsystem="nf_conntrack",
    name="insert_failed_vec",
    labels={"label_1", "label_2", "label_new"},
})
if not(err) then error("must be error") end


-- counter
local counter = prometheus.counter({
     namespace="node_scout",
     subsystem="nf_conntrack",
     name="insert_failed_counter",
     help="insert_failed from nf_conntrack",
})

counter:inc()
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed_counter 1")) then error("body:\n"..result.body) end

counter:add(2.2)
local result, err = client:do_request(request)
if err then error(err) end
if not(result.code == 200) then error("code:\n"..tostring(result.code)) end
if not(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed_counter 3.2")) then error("body:\n"..result.body) end
