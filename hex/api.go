// Package hex implements base64 encode/decode functionality for lua.
package hex

import (
	"encoding/hex"
	"io"

	lio "github.com/vadv/gopher-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
)

const (
	hexEncoderType = "hex.Encoder"
	hexDecoderType = "hex.Decoder"
)

// LVHexEncoder creates a new Lua user data for the hex encoder
func LVHexEncoder(L *lua.LState, writer io.Writer) lua.LValue {
	ud := L.NewUserData()
	ud.Value = writer
	L.SetMetatable(ud, L.GetTypeMetatable(hexEncoderType))
	return ud
}

// LVHexDecoder creates a new Lua user data for the hex decoder
func LVHexDecoder(L *lua.LState, reader io.Reader) lua.LValue {
	ud := L.NewUserData()
	ud.Value = reader
	L.SetMetatable(ud, L.GetTypeMetatable(hexDecoderType))
	return ud
}

// DecodeString decodes the encoded string with the encoding
func DecodeString(L *lua.LState) int {
	encoded := L.CheckString(1)
	L.Pop(L.GetTop())
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(decoded))
	return 1
}

// EncodeToString decodes the string with the encoding
func EncodeToString(L *lua.LState) int {
	decoded := L.CheckString(1)
	L.Pop(L.GetTop())
	encoded := hex.EncodeToString([]byte(decoded))
	L.Push(lua.LString(encoded))
	return 1
}

// registerHexEncoder Registers the encoder type and its methods
func registerHexEncoder(L *lua.LState) {
	mt := L.NewTypeMetatable(hexEncoderType)
	L.SetGlobal(hexEncoderType, mt)
	L.SetField(mt, "__index", lio.WriterFuncTable(L))
}

// registerHexDecoder Registers the decoder type and its methods
func registerHexDecoder(L *lua.LState) {
	mt := L.NewTypeMetatable(hexDecoderType)
	L.SetGlobal(hexDecoderType, mt)
	L.SetField(mt, "__index", lio.ReaderFuncTable(L))
}

func NewEncoder(L *lua.LState) int {
	writer := lio.CheckIOWriter(L, 1)
	L.Pop(L.GetTop())
	encoder := hex.NewEncoder(writer)
	L.Push(LVHexEncoder(L, encoder))
	return 1
}

func NewDecoder(L *lua.LState) int {
	reader := lio.CheckIOReader(L, 1)
	L.Pop(L.GetTop())
	decoder := hex.NewDecoder(reader)
	L.Push(LVHexDecoder(L, decoder))
	return 1
}
