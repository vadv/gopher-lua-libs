package time

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// time.sleep(number)
func ExampleSleep() {
	state := lua.NewState()
	Preload(state)
	source := `
    local time = require("time")
    local begin = time.unix()
    time.sleep(1.2)
    local stop = time.unix()
    local result = stop - begin
    -- round
    result = math.floor(result * 10^2 + 0.5) / 10^2
    print(result)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1
}

// time.parse(value, layout)
func ExampleParse() {
	state := lua.NewState()
	Preload(state)
	source := `
    local time = require("time")
    local result, err = time.parse("Dec  2 03:33:05 2018", "Jan  2 15:04:05 2006")
    if err then error(err) end
    print(result)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1543721585
}
