package crypto

import lua "github.com/yuin/gopher-lua"

// Preload adds crypto to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
// local crypto = require("crypto")
func Preload(L *lua.LState) {
	L.PreloadModule("crypto", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	t := L.NewTable()
	// Load the constants
	for name := range modeNames {
		t.RawSetString(name, lua.LString(name))
	}
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"md5":             MD5,
	"sha256":          SHA256,
	"aes_encrypt_hex": AESEncryptHex,
	"aes_decrypt_hex": AESDecryptHex,
	"aes_encrypt":     AESEncrypt,
	"aes_decrypt":     AESDecrypt,
}
