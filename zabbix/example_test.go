package zabbix

import (
	"log"

	http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

// example zabbix
func Example_package() {
	state := lua.NewState()
	Preload(state)
	http.Preload(state)
	inspect.Preload(state)
	source := `
local zabbix = require("zabbix")
local inspect = require("inspect")
local http = require("http")

local client = http.client({proxy="http://proxy"})
local zbx = zabbix.new({url="http://zabbix.url"}, client)


local err = zbx:login()
-- if err then error(err) end

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
-- if err then error(err) end

for k, v in pairs( response ) do
    print(inspect(v))
    print(v.description)
end
zbx:logout()
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}
