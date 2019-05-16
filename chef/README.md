# chef [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/chef?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/chef)

## Functions

- `client(client_name, path_to_file_with_key, chef_url, http_client_ud)` - returns chef client instance for further usage. Required [http](https://github.com/vadv/gopher-lua-libs/tree/master/http) client instance as `http_client_ud`. Please note that you must specify last slash in `chef_url`.

## Methods
### client
- `search(index, query, [partical_data], [params]` - make [search](https://docs.chef.io/api_chef_server.html#search) by given INDEX and query. Also possible use partical search and specify `offset`, `limit`, `order_by` in `params`.
- `request(verb, url)` - make request to chef server.


## Examples

```lua

local http = require("http")
local chef = require("chef")
local inspect = require("inspect")

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
print(inspect(result)) -- { [node_name_1=url_node_name_1], [node_name_2=url_node_name_2] }

-- get node
local result, err = client:request("GET", "nodes/node_name")
print(result.name)
-- chef node:
--  result.name
-- attributes
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
local result, err = client:search("node", "fqdn:node_name", {result_fqdn = {"fqdn"}})
print(inspect(result.rows)) -- table of results
print(result.rows[1]) -- {data = {result_fqdn = "node_name" }}

-- partical search node: limit = 2000, offset = 0, order by X_CHEF_id_CHEF_X asc
local result, err = client:search("node", "fqdn:node_name", {result_fqdn = {"fqdn"}}, {start=0, rows=2000})
```

