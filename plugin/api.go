// Package plugin implements run lua-code from lua-code.
package plugin

import (
	"context"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type luaPlugin struct {
	sync.Mutex
	*sync.Cond
	state      *lua.LState
	cancelFunc context.CancelFunc
	running    bool
	error      error
	body       *string
	filename   *string
	jobPayload *string
	args       []lua.LValue
	doneCh     lua.LChannel
	started    bool
}

func (p *luaPlugin) getError() error {
	p.Lock()
	defer p.Unlock()
	err := p.error
	return err
}

func (p *luaPlugin) getRunning() bool {
	p.Lock()
	defer p.Unlock()
	running := p.running
	return running
}

func (p *luaPlugin) setError(err error) {
	p.Lock()
	defer p.Unlock()
	p.error = err
}

func (p *luaPlugin) setRunning(val bool) {
	p.Lock()
	defer p.Unlock()
	p.running = val
}

// NewPluginState return lua state
func NewPluginState() *lua.LState {
	state := lua.NewState()
	PreloadAll(state)
	return state
}

func (p *luaPlugin) start() {
	p.Lock()
	state := NewPluginState()
	defer state.Close()
	p.state = state
	p.error = nil
	p.running = true
	p.doneCh = make(lua.LChannel, 1)
	p.started = true
	defer close(p.doneCh)
	isBody := p.filename == nil
	if p.jobPayload != nil {
		payload := *p.jobPayload
		state.SetGlobal(`payload`, lua.LString(payload))
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc
	p.state.SetContext(ctx)
	newArg := state.NewTable()
	for _, arg := range p.args {
		switch t := arg.Type(); t {
		case lua.LTFunction:
			arg = state.NewFunctionFromProto(arg.(*lua.LFunction).Proto)
		}
		newArg.Append(arg)
	}
	state.SetGlobal(`arg`, newArg)
	p.Signal()
	p.Unlock()

	// blocking
	if isBody {
		body := *p.body
		p.setError(p.state.DoString(body))
	} else {
		filename := *p.filename
		p.setError(p.state.DoFile(filename))
	}
	p.setRunning(false)
}

func checkPlugin(L *lua.LState, n int) *luaPlugin {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaPlugin); ok {
		return v
	}
	L.ArgError(n, "plugin expected")
	return nil
}

func NewLuaPlugin(L *lua.LState, n int) *luaPlugin {
	ret := &luaPlugin{}
	ret.Cond = sync.NewCond(&ret.Mutex)
	top := L.GetTop()
	for i := n; i <= top; i++ {
		arg := L.Get(i)
		switch t := arg.Type(); t {
		case lua.LTFunction:
			f := arg.(*lua.LFunction)
			if len(f.Upvalues) > 0 {
				L.ArgError(i, "cannot pass closures")
			}
			ret.args = append(ret.args, arg)
		case lua.LTTable:
			if L.GetMetatable(arg) != lua.LNil {
				L.ArgError(i, "tables with metadata are not allowed")
			}
			ret.args = append(ret.args, arg)
		case lua.LTNil, lua.LTBool, lua.LTNumber, lua.LTString, lua.LTChannel:
			ret.args = append(ret.args, arg)
		default:
			L.ArgError(i, t.String()+" is not allowed")
		}
	}
	return ret
}

// DoString lua plugin.do_string(body) returns plugin_ud
func DoString(L *lua.LState) int {
	body := L.CheckString(1)
	p := NewLuaPlugin(L, 2)
	p.body = &body
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable(`plugin_ud`))
	L.Push(ud)
	return 1
}

// DoFile lua plugin.do_file(filename) returns plugin_ud
func DoFile(L *lua.LState) int {
	filename := L.CheckString(1)
	p := NewLuaPlugin(L, 2)
	p.filename = &filename
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable(`plugin_ud`))
	L.Push(ud)
	return 1
}

// DoFileWithPayload lua plugin.async() returns (plugin_ud, err)
func DoFileWithPayload(L *lua.LState) int {
	filename := L.CheckString(1)
	payload := L.CheckString(2)
	p := NewLuaPlugin(L, 2)
	p.filename = &filename
	p.jobPayload = &payload
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable(`plugin_ud`))
	L.Push(ud)
	return 1
}

// DoStringWithPayload lua plugin.async() returns (plugin_ud, err)
func DoStringWithPayload(L *lua.LState) int {
	body := L.CheckString(1)
	payload := L.CheckString(2)
	p := NewLuaPlugin(L, 2)
	p.body = &body
	p.jobPayload = &payload
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable(`plugin_ud`))
	L.Push(ud)
	return 1
}

// Run lua plugin_ud:run()
func Run(L *lua.LState) (nRet int) {
	p := checkPlugin(L, 1)
	go p.start()

	// ensure it's started
	p.Lock()
	defer p.Unlock()
	for !p.started {
		p.Cond.Wait()
	}

	return
}

// IsRunning lua plugin_ud:is_running()
func IsRunning(L *lua.LState) int {
	p := checkPlugin(L, 1)
	L.Push(lua.LBool(p.getRunning()))
	return 1
}

// DoneChannel lua plugin_ud:done_chan()
func DoneChannel(L *lua.LState) int {
	p := checkPlugin(L, 1)
	if !p.started {
		L.ArgError(1, "Cannot obtain done channel on unstarted plugin")
	}
	L.Push(p.doneCh)
	return 1
}

// Error lua plugin_ud:error() returns err
func Error(L *lua.LState) int {
	p := checkPlugin(L, 1)
	err := p.getError()
	if err == nil {
		return 0
	}
	L.Push(lua.LString(err.Error()))
	return 1
}

// Stop lua plugin_ud:stop()
func Stop(L *lua.LState) int {
	p := checkPlugin(L, 1)
	p.Lock()
	defer p.Unlock()
	p.cancelFunc()
	return 0
}

// Wait lua plugin_ud:wait()
func Wait(L *lua.LState) int {
	p := checkPlugin(L, 1)
	// Can't wait if never started
	if !p.started {
		L.RaiseError("cannot wait on unstarted plugin")
	}
	// Add timeout if requested
	ctx := context.Background()
	cancel := func() {}
	timeout := lua.LVAsNumber(L.Get(2))
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout*lua.LNumber(time.Second)))
		defer cancel()
	}

	// Wait for done or timeout
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-p.doneCh:
		err = p.error
	}

	// return error
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
