package hex

import lua "github.com/yuin/gopher-lua"

// Preload adds hex to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//	local hex = require("hex")
func Preload(L *lua.LState) {
	L.PreloadModule("hex", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	registerHexDecoder(L)
	registerHexEncoder(L)

	// Register the encodings offered by hex go module.
	t := L.NewTable()
	L.SetFuncs(t, map[string]lua.LGFunction{
		"decode_string":    DecodeString,
		"encode_to_string": EncodeToString,
		"new_encoder":      NewEncoder,
		"new_decoder":      NewDecoder,
	})
	L.Push(t)
	return 1
}
