package hex

import (
	"log"

	"github.com/vadv/gopher-lua-libs/strings"
	lua "github.com/yuin/gopher-lua"
)

func ExampleEncodeToString() {
	state := lua.NewState()
	defer state.Close()
	Preload(state)
	source := `
    local hex = require 'hex'
	s = hex.encode_to_string("foo\01bar")
	print(s)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 666f6f01626172
}

func ExampleDecodeString() {
	state := lua.NewState()
	defer state.Close()
	Preload(state)
	source := `
    local hex = require 'hex'
	s, err = hex.decode_string("666f6f01626172")
	assert(not err, err)
	print(s)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// foobar
}

func ExampleNewEncoder() {
	state := lua.NewState()
	defer state.Close()
	Preload(state)
	strings.Preload(state)
	source := `
    local hex = require 'hex'
	local strings = require 'strings'
    local writer = strings.new_builder()
	encoder = hex.new_encoder(writer)
	encoder:write("foo\01bar")
	print(writer:string())
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// 666f6f01626172
}

func ExampleNewDecoder() {
	state := lua.NewState()
	defer state.Close()
	Preload(state)
	strings.Preload(state)
	source := `
    local hex = require 'hex'
	local strings = require 'strings'
    local reader = strings.new_reader('666f6f01626172')
	decoder = hex.new_decoder(reader)
    local s = decoder:read("*a")
	print(s)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// foobar
}
