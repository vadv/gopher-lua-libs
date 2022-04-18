package base64

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func ExampleEncodeToString() {
	state := lua.NewState()
	defer state.Close()
	Preload(state)
	source := `
    local base64 = require("base64")
	s = base64.RawStdEncoding:encode_to_string("foo\01bar")
	print(s)

	s = base64.StdEncoding:encode_to_string("foo\01bar")
	print(s)

	s = base64.RawURLEncoding:encode_to_string("this is a <tag> and should be encoded")
	print(s)

	s = base64.URLEncoding:encode_to_string("this is a <tag> and should be encoded")
	print(s)

`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// Zm9vAWJhcg
	// Zm9vAWJhcg==
	// dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA
	// dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==
}

func ExampleDecodeString() {
	state := lua.NewState()
	defer state.Close()
	Preload(state)
	source := `
    local base64 = require("base64")
	s, err = base64.RawStdEncoding:decode_string("Zm9vAWJhcg")
	assert(not err, err)
	print(s)

	s, err = base64.StdEncoding:decode_string("Zm9vAWJhcg==")
	assert(not err, err)
	print(s)

	s, err = base64.RawURLEncoding:decode_string("dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA")
	assert(not err, err)
	print(s)

	s, err = base64.URLEncoding:decode_string("dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==")
	assert(not err, err)
	print(s)

`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// foobar
	// foobar
	// this is a <tag> and should be encoded
	// this is a <tag> and should be encoded
}
