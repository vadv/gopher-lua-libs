// Package cloudwatch implements cloudwatch client api functionality for lua.
package cloudwatch

import (
	lua "github.com/yuin/gopher-lua"
)

// New lua new(profile, region) returns (clw_ud, err)
func New(L *lua.LState) int {
	var awsProfile, awsRegion *string
	if L.GetTop() > 0 {
		val := L.CheckString(1)
		awsProfile = &val
	}
	if L.GetTop() > 1 {
		val := L.CheckString(2)
		awsProfile = &val
	}
	clw, err := newLauClW(awsProfile, awsRegion)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = clw
	L.SetMetatable(ud, L.GetTypeMetatable("clw_ud"))
	L.Push(ud)
	return 1
}
