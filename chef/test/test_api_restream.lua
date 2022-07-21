local http = require("http")
local chef = require("chef")
local inspect = require("inspect")

local node = "chef.itv.restr.im"

local http_client = http.client({ insecure_ssl = true, timeout = 20 })
local client = chef.client(
        "lualibs",
        "./test/client.pem",
        "https://chef.itv.restr.im/organizations/restream/",
        http_client
)

function TestNodes(t)
    -- nodes
    local result, err = client:request("GET", "nodes")
    assert(not err, err)
    local found = false
    for k, v in pairs(result) do
        if k == node then
            found = true
        end
    end
    assert(found, "nodes don't contain " .. node)
end

function TestNodes(t)
    -- node
    local result, err = client:request("GET", "nodes/" .. node)
    assert(not err, err)
    t:Log(result.default)
    t:Log(result.automatic)
    t:Log(result.normal)
    t:Log(result.override)
    assert(result.name == node, "get node: " .. result.name)
end

function TestSearch(t)
    t:Run("1", function(t)
        local result, err = client:search("node", "fqdn:" .. node)
        assert(not err, err)
        assert(result.rows[1].name == node, "search node 1: " .. result.rows[1].name)
    end)

    t:Run("2", function(t)
        -- search 2
        local result, err = client:search("node", "fqdn:" .. node, { fqdn_name = { "name" } })
        assert(not err, err)
        t:Log(inspect(result.rows[1].data))
        assert(result.rows[1].data.fqdn_name == node, "search node 2: " .. tostring(result.rows[1].data.fqdn_name))
    end)
end
