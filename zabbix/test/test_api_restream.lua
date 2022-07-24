local zabbix = require("zabbix")
local inspect = require("inspect")
local http = require("http")

function Test_zabbix(t)
    local zbx = zabbix.new({ url = "http://192.168.184.83:8080", user = "zbgate", passwd = "EeTh9Aiw", debug = true })

    local err = zbx:login()
    assert(not err, err)

    --[[
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
    if response then
        for k, v in pairs(response) do
            if v.hosts and v.hosts[1] and v.items and v.items[1] and v.items[1].value_type == "3" then
                print(inspect(v))
                item_id = v.items[1].itemid
                print(v.hosts[1].host, v.description)
            end
        end
    end

    local err = zbx:save_graph(item_id, "./test/test.png")
    if err then error(err) end
    --]]

    local response, err = zbx:request("host.get",
            {
                filter = {
                    host = { "rt-core-pub01.core.iptv.rt.ru" }
                },
            }
    )
    assert(not err, err)
    assert(response, "must be result")

    local err = zbx:logout()
    assert(not err, err)
end
