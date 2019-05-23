package log

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// print(args..)
func Example_Print() {
	state := lua.NewState()
	Preload(state)
	source := `
    local log = require("log")
    local info = log.new("STDOUT", "[INFO] ")
    info:print("1 ", 2)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// [INFO] 1 2
}

// printf(string, args..)
func Example_Printf() {
	state := lua.NewState()
	Preload(state)
	source := `
    local log = require("log")
    local info = log.new("STDOUT", "[INFO] ")
    info:printf("%s %d\n", "1", 2)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// [INFO] 1 2
}

// println(string, args..)
func Example_Println() {
	state := lua.NewState()
	Preload(state)
	source := `
    local log = require("log")
    local info = log.new("STDOUT", "[INFO] ")
    info:println("1", 2)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// [INFO] 1 2
}

// set_flags(config={})
func Example_SetFlags() {
	state := lua.NewState()
	Preload(state)
	source := `
    local log = require("log")
    local logger = log.new()
    logger:set_prefix("[prefix] ")
    logger:set_flags({longfile=true})
    logger:println("1", 2)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// [prefix] <string>:6: 1 2
}
