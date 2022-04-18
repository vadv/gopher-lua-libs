// Package base64 implements base64 encode/decode functionality for lua.
package base64

import (
	"encoding/base64"
	lua "github.com/yuin/gopher-lua"
)

const (
	base64EncodingType = "base64.Encoding"
)

//CheckBase64Encoding checks the argument at position n is a *base64.Encoding
func CheckBase64Encoding(L *lua.LState, n int) *base64.Encoding {
	ud := L.CheckUserData(n)
	if encoding, ok := ud.Value.(*base64.Encoding); ok {
		return encoding
	}
	L.ArgError(n, base64EncodingType+" expected")
	return nil
}

//LVBase64Encoding converts encoding to a UserData type for lua
func LVBase64Encoding(L *lua.LState, encoding *base64.Encoding) lua.LValue {
	ud := L.NewUserData()
	ud.Value = encoding
	L.SetMetatable(ud, L.GetTypeMetatable(base64EncodingType))
	return ud
}

//DecodeString decodes the encoded string with the encoding
func DecodeString(L *lua.LState) int {
	encoding := CheckBase64Encoding(L, 1)
	encoded := L.CheckString(2)
	L.Pop(L.GetTop())
	decoded, err := encoding.DecodeString(encoded)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(decoded))
	return 1
}

//EncodeToString decodes the string with the encoding
func EncodeToString(L *lua.LState) int {
	encoding := CheckBase64Encoding(L, 1)
	decoded := L.CheckString(2)
	L.Pop(L.GetTop())
	encoded := encoding.EncodeToString([]byte(decoded))
	L.Push(lua.LString(encoded))
	return 1
}

//registerBase64Encoding Registers the encoding type and its methods
func registerBase64Encoding(L *lua.LState) {
	mt := L.NewTypeMetatable(base64EncodingType)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"decode_string":    DecodeString,
		"encode_to_string": EncodeToString,
	}))
}
