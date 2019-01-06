# chef [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/chef?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/chef)

## Usage

```lua

local http = require("http")
local chef = require("chef")

local http_client = http.client({insecure_ssl=true})
local client = chef.client(
    "client_name",
    "path/to/client.pem",
    "https://chef.org/organizations/org/",
    http_client
)

-- list nodes
local result, err = client:request("GET", "nodes")
if err then error(err) end
print(result[1][1]) -- { {node_name_1=url_node_name_1}, {node_name_2=url_node_name_2} }

-- get node
local result, err = client:request("GET", "nodes/node_name")
print(result.name)
-- chef node:
-- result.name
--  attributes
--  result.automatic
--  result.default
--  result.normal
--  result.override


-- search node
local result, err = client:search("node", "fqdn:node_name")
print(result.total) -- total count
print(result.start) -- offset
print(inspect(result.rows)) -- table of nodes

-- partical search node
local result, err = client:search("node", "fqdn:node_name", {"result_fqdn" = {"fqdn"}})
print(inspect(result.rows)) -- table of results
print(result.rows[0]) -- {data = {result_fqdn = "node_name" }}

-- partical search node: limit = 2000, offset = 0, order by X_CHEF_id_CHEF_X asc
local result, err = client:search("node", "fqdn:node_name", {"result_fqdn" = {"fqdn"}}, {start=0, rows=2000})
```

