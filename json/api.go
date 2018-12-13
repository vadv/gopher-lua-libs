// Package json implements json decode/encode functionality for lua.
package json

import (
	"encoding/json"
	"errors"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

var (
	ErrJsonFunction = errors.New("cannot encode function to JSON")
	ErrJsonChannel  = errors.New("cannot encode channel to JSON")
	ErrJsonState    = errors.New("cannot encode state to JSON")
	ErrJsonUserData = errors.New("cannot encode userdata to JSON")
	ErrJsonNested   = errors.New("cannot encode recursively nested tables to JSON")
)

type luaJson struct {
	lua.LValue
	visited map[*lua.LTable]bool
}

func (j luaJson) MarshalJSON() ([]byte, error) {
	return toJSON(j.LValue, j.visited)
}

// Decode(): lua json.decode(string) returns (table, err)
func Decode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := json.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(fromJSON(L, value))
	return 1
}

// Encode(): lua json.encode(obj) returns (string, err)
func Encode(L *lua.LState) int {
	value := L.CheckAny(1)

	visited := make(map[*lua.LTable]bool)
	data, err := toJSON(value, visited)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(string(data)))
	return 1
}

func toJSON(value lua.LValue, visited map[*lua.LTable]bool) (data []byte, err error) {
	switch converted := value.(type) {
	case lua.LBool:
		data, err = json.Marshal(converted)
	case lua.LChannel:
		err = ErrJsonChannel
	case lua.LNumber:
		data, err = json.Marshal(converted)
	case *lua.LFunction:
		err = ErrJsonFunction
	case *lua.LNilType:
		data, err = json.Marshal(converted)
	case *lua.LState:
		err = ErrJsonState
	case lua.LString:
		data, err = json.Marshal(converted)
	case *lua.LTable:
		var arr []luaJson
		var obj map[string]luaJson

		if visited[converted] {
			panic(ErrJsonNested)
		}
		visited[converted] = true

		converted.ForEach(func(k lua.LValue, v lua.LValue) {
			i, numberKey := k.(lua.LNumber)
			if numberKey && obj == nil {
				index := int(i) - 1
				if index != len(arr) {
					// map out of order; convert to map
					obj = make(map[string]luaJson)
					for i, value := range arr {
						obj[strconv.Itoa(i+1)] = value
					}
					obj[strconv.Itoa(index+1)] = luaJson{v, visited}
					return
				}
				arr = append(arr, luaJson{v, visited})
				return
			}
			if obj == nil {
				obj = make(map[string]luaJson)
				for i, value := range arr {
					obj[strconv.Itoa(i+1)] = value
				}
			}
			obj[k.String()] = luaJson{v, visited}
		})
		if obj != nil {
			data, err = json.Marshal(obj)
		} else {
			data, err = json.Marshal(arr)
		}
	case *lua.LUserData:
		// TODO: call metatable __tostring?
		err = ErrJsonUserData
	}
	return
}

func fromJSON(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromJSON(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(lua.LString(key), fromJSON(L, item))
		}
		return tbl
	}
	return lua.LNil
}
