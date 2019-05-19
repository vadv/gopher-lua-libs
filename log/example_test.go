package log

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func Example_package() {
	state := lua.NewState()
	Preload(state)
	source := `
local log = require("log")
local logger = log.logger({flags=100000, level='debug'})
logger:debug('2nd DEBUG')
logger:info('2nd INFO')
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// [DEBUG] 2nd DEBUG
	// [INFO] 2nd INFO
}
