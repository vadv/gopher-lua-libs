package json

import (
	"log"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

// json.decode(string)
func ExampleDecode() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	source := `
    local json = require("json")
    local inspect = require("inspect")
    local jsonString = [[{"a":{"b":1}}]]
    local result, err = json.decode(jsonString)
    if err then error(err) end
	local remove_all_metatables = function(item, path)
	  if path[#path] ~= inspect.METATABLE then return item end
	end
    print(inspect(result, {process = remove_all_metatables, newline="", indent=""}))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// {a = {b = 1}}
}

// json.encode(obj)
func ExampleEncode() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	source := `
    local json = require("json")
    local inspect = require("inspect")
    local table = {a={b=1}}
    local result, err = json.encode(table)
    if err then error(err) end
    print(inspect(result, {newline="", indent=""}))

	print(inspect( json.encode( {} ) ))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// '{"a":{"b":1}}'
	// "[]"
}
