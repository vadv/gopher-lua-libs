// Package cloudwatch implements cloudwatch client api functionality for lua.
package cloudwatch

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
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

// Download lua clw:download(filename, filter={}, timeout) error
// filter table:
//   {
//     log_group_name="/aws/rds/instance/name/rds",
//     filter_patern="",
//     start_time=timestamp,
//     end_time=timestamp,
//   }
func Download(L *lua.LState) int {

	clw := checkluaClW(L, 1)

	// parse paramters
	filename := L.CheckString(2)
	filterLua := L.CheckTable(3)
	filter := &cloudwatchlogs.FilterLogEventsInput{}
	filterLua.ForEach(func(k lua.LValue, v lua.LValue) {
		if k.String() == `log_group_name` {
			filter.SetLogGroupName(v.String())
		}
		if k.String() == `filter_patern` {
			filter.SetFilterPattern((v.String())
		}
		if k.String() == `start_time` {
			if value, ok := v.(lua.LNumber); ok {
				setting := int64(float64(value) * 1000)
				filter.SetStartTime(setting)
			} else {
				L.ArgError(3, "start must be number")
			}
		}
		if k.String() == `end_time` {
			if value, ok := v.(lua.LNumber); ok {
				setting := int64(float64(value) * 1000)
				filter.SetEndTime(setting)
			} else {
				L.ArgError(3, "start must be number")
			}
		}
	})

	timeout := 60
	if L.GetTop() > 5 {
		timeout = L.CheckInt(6)
	}

	events := make(chan *cloudwatchlogs.FilteredLogEvent, 10)
	done := make(chan error, 1)

	fd, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	defer fd.Close()

	go clw.events(filter, events, done)

	for {
		select {
		case event := <-events:
			if event != nil {
				msg := event.Message
				if msg != nil {
					if _, err := fd.WriteString(*msg + "\n"); err != nil {
						L.Push(lua.LString(err.Error()))
						return 1
					}
				}
			}
		case err := <-done:
			if err != nil {
				L.Push(lua.LString(err.Error()))
				return 1
			}
			return 0
		case <-time.After(time.Duration(timeout) * time.Second):
			L.Push(lua.LString("timeout"))
			return 1
		}
	}
}
