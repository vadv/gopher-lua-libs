package base64

import (
	"encoding/base64"
	lua "github.com/yuin/gopher-lua"
)

// Preload adds yaml to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local yaml = require("yaml")
func Preload(L *lua.LState) {
	L.PreloadModule("base64", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	registerBase64Encoding(L)

	// Register the encodings offered by base64 go module.
	t := L.NewTable()
	L.SetField(t, "RawStdEncoding", LVBase64Encoding(L, base64.RawStdEncoding))
	L.SetField(t, "RawURLEncoding", LVBase64Encoding(L, base64.RawURLEncoding))
	L.SetField(t, "StdEncoding", LVBase64Encoding(L, base64.StdEncoding))
	L.SetField(t, "URLEncoding", LVBase64Encoding(L, base64.URLEncoding))

	// TODO(scr): When https://github.com/vadv/gopher-lua-libs/pull/29 lands, Add NewEncoder/Decoder methods so that
	// 			  encoding/decoding can be done directly from/to files.
	L.Push(t)
	return 1
}
