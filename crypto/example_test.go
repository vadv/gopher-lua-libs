package crypto

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// crypto.md5(string)
func ExampleMD5() {
	state := lua.NewState()
	Preload(state)
	source := `
    local crypto = require("crypto")
    print(crypto.md5("1\n"))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// b026324c6904b2a9cb4b88d6d61c81d1
}

// crypto.sha256(string)
func ExampleSHA256() {
	state := lua.NewState()
	Preload(state)
	source := `
    local crypto = require("crypto")
    print(crypto.sha256("1\n"))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 4355a46b19d348dc2f57c046f8ef63d4538ebb936000f3c9ee954a27460dd865
}
