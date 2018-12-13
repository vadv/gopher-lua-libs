package strings

import (
	"log"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

// strings.split(string, sep)
func ExampleSplit() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	source := `
    local inspect = require("inspect")
    local strings = require("strings")
    local result = strings.split("a b c d", " ")
    print(inspect(result, {newline="", indent=""}))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// { "a", "b", "c", "d" }
}

// strings.has_prefix(string, prefix)
func ExampleHasPrefix() {
	state := lua.NewState()
	Preload(state)
	source := `
    local strings = require("strings")
    local result = strings.has_prefix("abcd", "a")
    print(result)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// true
}

// strings.has_suffix(string, suffix)
func ExampleHasSuffix() {
	state := lua.NewState()
	Preload(state)
	source := `
    local strings = require("strings")
    local result = strings.has_suffix("abcd", "d")
    print(result)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// true
}

// strings.trim(string, cutset)
func ExampleTrim() {
	state := lua.NewState()
	Preload(state)
	source := `
    local strings = require("strings")
    local result = strings.trim("abcd", "d")
    print(result)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// abc
}
