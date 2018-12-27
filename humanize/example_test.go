package humanize

import (
	"log"

	time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

// humanize.ibytes(number)
func ExampleIBytes() {
	state := lua.NewState()
	Preload(state)
	source := `
local humanize = require("humanize")
print(humanize.ibytes(1395864371))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1.3 GiB
}

// humanize.parse_bytes(string)
func ExampleParseBytes() {
	state := lua.NewState()
	Preload(state)
	source := `
local humanize = require("humanize")
print(humanize.parse_bytes("1.3GiB"))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1395864371
}

// humanize.time(number)
func ExampleTime() {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	source := `
local humanize = require("humanize")
local time = require("time")
print(humanize.time(time.unix() - 61))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1 minute ago
}
