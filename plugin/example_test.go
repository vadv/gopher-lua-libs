package plugin

import (
	"log"

	time "github.com/vadv/gopher-lua-libs/time"

	lua "github.com/yuin/gopher-lua"
)

// plugin.do_string(), plugin_ud:run(), plugin_ud:stop()
func Example_package() {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	source := `
    local plugin = require("plugin")
    local time = require("time")

    local plugin_body = [[
        local time = require("time")
        local i = 1
        while true do
            print(i)
            i = i + 1
            time.sleep(1)
        end
    ]]

    local print_plugin = plugin.do_string(plugin_body)
    print_plugin:run()
    time.sleep(2)
    print_plugin:stop()
    time.sleep(1)

    local running = print_plugin:is_running()
    if running then error("already running") end

`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1
	// 2
}
