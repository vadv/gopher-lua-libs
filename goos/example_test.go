package goos

import (
	"log"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

// goos.stat(filename)
func ExampleStat() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	source := `
local goos = require("goos")
local inspect = require("inspect")

local info, err = goos.stat("./test/test.file")
if err then error(err) end
info.mode=""
print(inspect(info, {newline="", indent=""}))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// {is_dir = false,mod_time = 1545917050,mode = "",size = 0}
}
