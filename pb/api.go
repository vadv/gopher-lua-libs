// Package pb implements github.com/cheggaaa/pb/v3 functionality for lua.
package pb

import (
	lua "github.com/yuin/gopher-lua"

	"fmt"
	"github.com/cheggaaa/pb/v3"
	"os"
	"strconv"
	"time"
)

// New (count) - total for progress bar
func New(L *lua.LState) int {
	var bar *pb.ProgressBar

	if L.GetTop() > 0 {
		bar = pb.New(L.CheckInt(1))
	}

	ud := L.NewUserData()
	ud.Value = bar
	L.SetMetatable(ud, L.GetTypeMetatable(`pb_ud`))
	L.Push(ud)
	return 1
}

func checkBar(L *lua.LState, n int) *pb.ProgressBar {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*pb.ProgressBar); ok {
		return v
	}
	L.ArgError(n, "progress bar expected")
	return nil
}

func Configure(L *lua.LState) int {
	bar := checkBar(L, 1)
	config := L.CheckTable(2)
	var err error
	config.ForEach(func(k lua.LValue, v lua.LValue) {
		switch k.String() {
		case `template`:
			var tmpl string
			switch v.String() {
			case `full`:
				tmpl = string(pb.Full)
				bar.SetTemplateString(tmpl)
			case `simple`:
				tmpl = string(pb.Simple)
				bar.SetTemplateString(tmpl)

			default:
				// Custom template
				bar.SetTemplateString(v.String())
				if err = bar.Err(); err != nil {
					err = fmt.Errorf("error while setting template: `%s`", err)
				}
			}
		case `writer`:
			switch v.String() {
			case `stdout`:
				bar.SetWriter(os.Stdout)
			default:
				err = fmt.Errorf("unknown writer: `%s`", v.String())
			}
		case `refresh_rate`:
			if v.Type() != lua.LTNumber {
				err = fmt.Errorf(`refresh_rate must be an integer`)
				L.Push(lua.LString(err.Error()))
				return
			}
			i, _ := strconv.ParseInt(v.String(), 10, 64)
			duration := time.Duration(i)
			bar.SetRefreshRate(duration * time.Millisecond)

		default:
			err = fmt.Errorf("unknown config parameter: `%s`", k.String())
		}
	})

	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

func Start(L *lua.LState) int {
	bar := checkBar(L, 1)
	bar.Start()
	return 1
}

func Increment(L *lua.LState) int {
	bar := checkBar(L, 1)
	bar.Increment()
	return 1
}

func Finish(L *lua.LState) int {
	bar := checkBar(L, 1)
	bar.Finish()
	return 1
}
