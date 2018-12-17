package storage

import (
	"log"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	time "github.com/vadv/gopher-lua-libs/time"

	lua "github.com/yuin/gopher-lua"
)

// plugin.new(), plugin_ud:run(), plugin_ud:stop()
func Example_package() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	time.Preload(state)
	source := `
local storage = require("storage")
local inspect = require("inspect")

local s, err = storage.open("./test/db.json")
if err then error(err) end

local err = s:set("key", {"one", "two", 1}, 10)
if err then error(err) end

local value, found, err = s:get("key")
if err then error(err) end
if not found then error("must be found") end

print(  inspect(value, {newline="", indent=""}) )


`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// { "one", "two", 1 }
}
