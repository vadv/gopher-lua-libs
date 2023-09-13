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
	// 1.2
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

// time.format(value, layout, timezone)
func ExampleFormat() {
	state := lua.NewState()
	Preload(state)
	source := `
    local time = require("time")
    print( time.format(0, "Mon Jan 2 15:04:05 -0700 MST 2006", "UTC") )
    print( time.format(0, "Mon Jan 2 15:04:05 -0700 MST 2006", "Europe/Moscow") )
    print( time.format(1543721585, "Jan  2 15:04:05 2006", "Europe/Moscow") )
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// Thu Jan 1 00:00:00 +0000 UTC 1970
	// Thu Jan 1 03:00:00 +0300 MSK 1970
	// Dec  2 06:33:05 2018
}
