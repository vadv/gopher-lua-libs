local prometheus = require("prometheus")
local http = require("http_client")
local strings = require("strings")
local time = require("time")

function Test_prometheus(t)
    local pp = prometheus.register(":18080")
    pp:start()

    local client = http.client({ timeout = 5 })
    local request = http.request("GET", "http://127.0.0.1:18080/metrics")

    time.sleep(1)

    t:Run("gauge", function(t)
        local gauge = prometheus.gauge({
            namespace = "node_scout",
            subsystem = "nf_conntrack",
            name = "insert_failed",
            help = "insert_failed from nf_conntrack",
        })
        gauge:set(100)

        local result, err = client:do_request(request)
        assert(not err, err)
        assert(result.code == 200, "code:\n" .. tostring(result.code))
        assert(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed 100", "body:\n" .. result.body))

        gauge:add(1)
        local result, err = client:do_request(request)
        assert(not err, err)
        assert(result.code == 200, "code:\n" .. tostring(result.code))
        assert(strings.contains(result.body, "node_scout_nf_conntrack_insert_failed 101", "body:\n" .. result.body))
    end)

    t:Run("gauge vector", function(t)
        -- gauge vector
        local gauge_vec, err = prometheus.gauge({
            namespace = "node_scout",
            subsystem = "nf_conntrack",
            name = "insert_failed_vec",
            help = "insert_failed from nf_conntrack",
            labels = { "label_1", "label_2" },
        })
        assert(not err, err)

        gauge_vec:set(100, { label_1 = "one", label_2 = "two" })
        local result, err = client:do_request(request)
        assert(not err, err)
        assert(result.code == 200, "code:\n" .. tostring(result.code))
        assert(strings.contains(result.body,
                'node_scout_nf_conntrack_insert_failed_vec{label_1="one",label_2="two"} 100', "body:\n" .. result.body))

        gauge_vec:add(1, { label_1 = "one", label_2 = "two" })
        local result, err = client:do_request(request)
        assert(not err, err)
        assert(result.code == 200, "code:\n" .. tostring(result.code))
        assert(strings.contains(result.body,
                'node_scout_nf_conntrack_insert_failed_vec{label_1="one",label_2="two"} 101', "body:\n" .. result.body))

        gauge_vec:inc({ label_1 = "one", label_2 = "two" })
        local result, err = client:do_request(request)
        assert(not err, err)
        assert(result.code == 200, "code:\n" .. tostring(result.code))
        assert(strings.contains(result.body,
                'node_scout_nf_conntrack_insert_failed_vec{label_1="one",label_2="two"} 102', "body:\n" .. result.body))
    end)

    t:Run("re-register gauge vector", function(t)
        local _, err = prometheus.gauge({
            namespace = "node_scout",
            subsystem = "nf_conntrack",
            name = "insert_failed_vec",
            labels = { "label_1", "label_2", "label_new" },
        })
        assert(err, "must be error")
    end)

    t:Run("counter", function(t)
        -- counter
        local counter = prometheus.counter({
            namespace = "node_scout",
            subsystem = "nf_conntrack",
            name = "insert_failed_counter",
            help = "insert_failed from nf_conntrack",
        })

        counter:inc()
        local result, err = client:do_request(request)
        assert(not err, err)
        assert(result.code == 200, "code:\n" .. tostring(result.code))
        assert(strings.contains(result.body,
                "node_scout_nf_conntrack_insert_failed_counter 1", "body:\n" .. result.body))

        counter:add(2.2)
        local result, err = client:do_request(request)
        assert(not err, err)
        assert(result.code == 200, "code:\n" .. tostring(result.code))
        assert(strings.contains(result.body,
                "node_scout_nf_conntrack_insert_failed_counter 3.2", "body:\n" .. result.body))
    end)
end
