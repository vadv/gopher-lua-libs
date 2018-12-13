package cmd

import (
	"log"

	runtime "github.com/vadv/gopher-lua-libs/runtime"
	lua "github.com/yuin/gopher-lua"
)

// cmd.exec()
func ExampleExec() {
	state := lua.NewState()
	Preload(state)
	runtime.Preload(state)
	source := `
local cmd = require("cmd")
local runtime = require("runtime")

local command = "sleep 1"
if runtime.goos() == "windows" then command = "timeout 1" end

local result, err = cmd.exec(command)
if err then error(err) end
print(result.status)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 0
}
