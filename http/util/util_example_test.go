package util

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// http_clien.parse_url(string)
func ExampleParseUrl() {
	state := lua.NewState()
	Preload(state)
	source := `
    local http_util = require("http_util")
    local url, err = http_util.parse_url("http://u1:p2@host:80/pathx?k1=v1&k2=v2&k1=vx")
    if err then error(err) end
    print(url.path)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// /pathx
}

// http_clien.build_url(table)
func ExampleBuidUrl() {
	state := lua.NewState()
	Preload(state)
	source := `
    local http_util = require("http_util")
    local url, err = http_util.parse_url("http://u1:p2@host:80/pathx?k1=v1&k2=v2&k1=vx")
    if err then error(err) end
    url.path = "path2"
    print(http_util.build_url(url))
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// http://u1:p2@host:80/path2?k1=v1&k1=vx&k2=v2
}
