// Package yaml implements yaml decode functionality for lua.
package yaml

import (
	lua "github.com/yuin/gopher-lua"
	yaml "gopkg.in/yaml.v2"
)

// Decode lua yaml.decode(string) returns (table, error)
func Decode(L *lua.LState) int {
	str := L.CheckString(1)

	var value interface{}
	err := yaml.Unmarshal([]byte(str), &value)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(fromYAML(L, value))
	return 1
}

func fromYAML(L *lua.LState, value interface{}) lua.LValue {
	switch converted := value.(type) {
	case bool:
		return lua.LBool(converted)
	case float64:
		return lua.LNumber(converted)
	case int:
		return lua.LNumber(converted)
	case int64:
		return lua.LNumber(converted)
	case string:
		return lua.LString(converted)
	case []interface{}:
		arr := L.CreateTable(len(converted), 0)
		for _, item := range converted {
			arr.Append(fromYAML(L, item))
		}
		return arr
	case map[interface{}]interface{}:
		tbl := L.CreateTable(0, len(converted))
		for key, item := range converted {
			tbl.RawSetH(fromYAML(L, key), fromYAML(L, item))
		}
		return tbl
	case interface{}:
		if v, ok := converted.(bool); ok {
			return lua.LBool(v)
		}
		if v, ok := converted.(float64); ok {
			return lua.LNumber(v)
		}
		if v, ok := converted.(string); ok {
			return lua.LString(v)
		}
	}
	return lua.LNil
}
