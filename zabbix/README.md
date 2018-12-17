# zabbix [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/zabbix?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/zabbix)

## Usage

```lua
local zabbix = require("zabbix")
local inspect = require("inspect")
local http = require("http")

local client = http.client({proxy="http://proxy", insecure_ssl=true, basic_auth_user="", basic_auth_password=""}) -- override default http client
local zbx = zabbix.new({url="https://zabbix.url"}, client)

local err = zbx:login()
if err then error(err) end

local response, err = zbx:request("trigger.get",
    {
        selectHosts = "extend", selectItems = "extend", selectLastEvent="extend",
        output = "extend", sortfield = "priority",
        filter = {
            sortorder="DESC", value="1", status=0
        },
        expandData = "1"
    }
)
if err then error(err) end

local item_id = 0
for k, v in pairs(response) do
    if v.hosts and v.hosts[1] and v.items and v.items[1] and v.items[1].value_type == "3" then
        -- print(inspect(v))
        item_id = v.items[1].itemid
        print(v.hosts[1].host, v.description)
    end
end

local err = zbx:save_graph(item_id, "./test/test.png")
if err then error(err) end
```

