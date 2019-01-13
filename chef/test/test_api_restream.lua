local http = require("http")
local chef = require("chef")
local inspect = require("inspect")

local node = "chef.itv.restr.im"

local http_client = http.client({insecure_ssl=true, timeout=20})
local client = chef.client(
    "lualibs",
    "./test/client.pem",
    "https://chef.itv.restr.im/organizations/restream/",
    http_client
)

-- nodes
local result, err = client:request("GET", "nodes")
if err then error(err) end
local found = false
for k,v in pairs(result) do if k == node then found = true end end
if not found then error("nodes") end
print("done: nodes")

-- node
local result, err = client:request("GET", "nodes/"..node)
if err then error(err) end
if not (result.name == node) then error("get node") end
print(result.default)
print(result.automatic)
print(result.normal)
print(result.override)
print("done: get node")

-- search 1
local result, err = client:search("node", "fqdn:"..node)
if err then error(err) end
if not (result.rows[1].name == node) then error("search node 1") end
print("done: search 1")

-- search 2
local result, err = client:search("node", "fqdn:"..node, {fqdn_name={"name"}})
if err then error(err) end
print(inspect(result.rows[1].data))
if not (result.rows[1].data.fqdn_name == node) then error("search node 2: "..tostring(result.rows[1].data.fqdn_name)) end
print("done: search 2")

