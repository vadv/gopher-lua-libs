package template

import (
	"log"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

func Example_package() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	source := `
local template = require("template")

local mustache, err = template.choose("mustache")

local values = {name="world"}
print( mustache:render("Hello {{name}}!", values) ) -- mustache:render_file(filename values)

local values = {data = {"one", "two"}}
print( mustache:render("{{#data}} {{.}} {{/data}}", values) )
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// Hello world!
	//  one  two
}
