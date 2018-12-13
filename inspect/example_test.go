package inspect

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// inspect(obj)
func Example_full() {
	state := lua.NewState()
	Preload(state)
	source := `
	local inspect = require("inspect")
    local table = {a={b=2}}
    print(inspect(table, {newline="", indent=""}))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	//{a = {b = 2}}
}
