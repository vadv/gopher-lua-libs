package cert_util

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// cert_util.not_after("host", <ip:port>)
func Example_package() {
	state := lua.NewState()
	Preload(state)
	source := `
    local cert_util = require("cert_util")
    local tx, err = cert_util.not_after("google.com", "64.233.165.101:443")
    if err then error(err) end
    print(tx > 0)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// true
}
