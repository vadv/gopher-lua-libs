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
info.mod_time=0
print(inspect(info, {newline="", indent=""}))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// {is_dir = false,mod_time = 0,mode = "",size = 0}
}

// goos.hostname()
func ExampleHostname() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	source := `
local goos = require("goos")
local hostname, err = goos.hostname()
if err then error(err) end
print(hostname > "")
	`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// true
}
