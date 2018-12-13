package tac

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// tac.open(), tac_ud:line(), tac_ud:close()
func Example() {
	state := lua.NewState()
	Preload(state)
	source := `
        local file = io.open("./test/file.txt", "w")
        file:write("1", "\n")
        file:write("2", "\n")
        file:write("3", "\n")

        local tac = require("tac")
        local scanner, err = tac.open("./test/file.txt")
        if err then error(err) end

        while true do
            local line = scanner:line()
            if line == nil then break end
            print(line)
        end
        scanner:close()
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 3
	// 2
	// 1
}
