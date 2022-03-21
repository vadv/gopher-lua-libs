package shellescape

import (
	lua "github.com/yuin/gopher-lua"
	"log"
)

func ExampleQuote() {
	L := lua.NewState()
	Preload(L)
	source := `
local shellescape = require("shellescape")
print(shellescape.quote("foo"))
`
	if err := L.DoString(source); err != nil {
		log.Fatal(err)
	}
	// Output:
	// foo
}

func ExampleQuoteCommand() {
	L := lua.NewState()
	Preload(L)
	source := `
local shellescape = require("shellescape")
print(shellescape.quote_command({"echo", "foo bar baz"}))
`
	if err := L.DoString(source); err != nil {
		log.Fatal(err)
	}
	// Output:
	// echo 'foo bar baz'
}

func ExampleStripUnsafe() {
	L := lua.NewState()
	Preload(L)
	source := `
local shellescape = require("shellescape")
print(shellescape.strip_unsafe("foo\nbar"))
`
	if err := L.DoString(source); err != nil {
		log.Fatal(err)
	}
	// Output:
	// foobar
}
