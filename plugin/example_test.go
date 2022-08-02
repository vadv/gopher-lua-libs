package plugin

import (
	time "github.com/vadv/gopher-lua-libs/time"
	"log"

	lua "github.com/yuin/gopher-lua"
)

// plugin.do_string(), plugin_ud:run(), plugin_ud:stop()
func Example_package() {
	state := lua.NewState()
	defer state.Close()
	Preload(state)
	time.Preload(state)
	source := `
    local plugin = require("plugin")
    local time = require("time")

    local plugin_body = [[
        local doCh, doneCh = unpack(arg)
        local i = 1
        while doCh:receive() do
            print(i)
            doneCh:send(i)
            i = i + 1
        end
    ]]

    -- Make synchronization channels and fire up the plugin
    local doCh = channel.make(100)
    local doneCh = channel.make(100)
    local print_plugin = plugin.do_string(plugin_body, doCh, doneCh)
    print_plugin:run()

    -- Allow two iterations to proceed
    doCh:send(nil)
    local ok, got = doneCh:receive()
    assert(ok and got == 1, string.format("ok = %s; got = %s", ok, got))
    doCh:send(nil)
    ok, got = doneCh:receive()
    assert(ok and got == 2, string.format("ok = %s; got = %s", ok, got))

    -- Close the doCh and wait to ensure it's closed gracefully but stop just to be sure
    doCh:close()
    time.sleep(1)
    print_plugin:stop()
    time.sleep(1)

    -- Ensure it's not still running
    assert(not print_plugin:is_running(), "still running")
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1
	// 2
}

func Example_using_like_goroutine() {
	state := lua.NewState()
	defer state.Close()

	Preload(state)

	source := `
	local plugin = require 'plugin'

	function myfunc(...)
		print(table.concat({...}, ""))
	end

	local background_body = [[
		return pcall(unpack(arg))
	]]
	local background_plugin = plugin.do_string(background_body, myfunc, "Hello", " ", "World")
	background_plugin:run()
	local err = background_plugin:wait()
	assert(not err, err)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err)
	}
	// Output:
	// Hello World
}
