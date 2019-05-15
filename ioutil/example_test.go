package ioutil

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// ioutil.read_file(filepath)
func ExampleReadFile() {
	state := lua.NewState()
	Preload(state)
	source := `
    local file = io.open("./test/file.data", "w")
    file:write("content of test file", "\n")
    file:close()


    local ioutil = require("ioutil")
    local result, err = ioutil.read_file("./test/file.data")
    if err then error(err) end
    print(result)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// content of test file
}

// ioutil.write_file(filepath)
func ExampleWriteFile() {
	state := lua.NewState()
	Preload(state)
	source := `
    local ioutil = require("ioutil")
    local err = ioutil.write_file("./test/file.data", "content of test file")
    if err then error(err) end
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	//
}
