package json

import (
	"encoding/json"
	"github.com/vadv/gopher-lua-libs/io"
	lua "github.com/yuin/gopher-lua"
)

const (
	jsonDecoderType = "json.Decoder"
)

func CheckDecoder(L *lua.LState, n int) *json.Decoder {
	ud := L.CheckUserData(n)
	if decoder, ok := ud.Value.(*json.Decoder); ok {
		return decoder
	}
	L.ArgError(n, jsonDecoderType+" expected")
	return nil
}

func LVDecoder(L *lua.LState, decoder *json.Decoder) lua.LValue {
	ud := L.NewUserData()
	ud.Value = decoder
	L.SetMetatable(ud, L.GetTypeMetatable(jsonDecoderType))
	return ud
}

func DecoderDecode(L *lua.LState) int {
	decoder := CheckDecoder(L, 1)
	L.Pop(L.GetTop())
	var value interface{}
	if err := decoder.Decode(&value); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(decode(L, value))
	return 1
}

func DecoderInputOffset(L *lua.LState) int {
	decoder := CheckDecoder(L, 1)
	L.Pop(L.GetTop())
	L.Push(lua.LNumber(decoder.InputOffset()))
	return 1
}

func DecoderMore(L *lua.LState) int {
	decoder := CheckDecoder(L, 1)
	L.Pop(L.GetTop())
	L.Push(lua.LBool(decoder.More()))
	return 1
}

func registerDecoder(L *lua.LState) {
	mt := L.NewTypeMetatable(jsonDecoderType)
	L.SetGlobal(jsonDecoderType, mt)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"decode":       DecoderDecode,
		"input_offset": DecoderInputOffset,
		"more":         DecoderMore,
	}))
}

func NewDecoder(L *lua.LState) int {
	reader := io.CheckReader(L, 1)
	L.Pop(L.GetTop())
	decoder := json.NewDecoder(reader)
	L.Push(LVDecoder(L, decoder))
	return 1
}
