package xmlpath

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

// xmlpath.compile, xmlpath.load, xmlpath_node_ud, xmlpath_path_ud, xmlpath_iter_ud
func Example_full() {
	state := lua.NewState()
	Preload(state)
	source := `
local xmlpath = require("xmlpath")

local data = [[
<channels>
    <channel id="1" xz1="600" />
    <channel id="2"           />
    <channel id="x" xz2="600" />
</channels>
]]
local data_path = "//channel/@id"

local node, err = xmlpath.load(data)
if err then error(err) end

local path, err = xmlpath.compile(data_path)
if err then error(err) end

local iter = path:iter(node)

for k, v in pairs(iter) do print(v:string()) end
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 1
	// 2
	// x
}
