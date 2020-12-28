// Package plugin implements run lua-code from lua-code.
package plugin

import (
	"context"
	"sync"

	"github.com/vadv/gopher-lua-libs/stats"

	cloudwatch "github.com/vadv/gopher-lua-libs/aws/cloudwatch"
	chef "github.com/vadv/gopher-lua-libs/chef"
	cmd "github.com/vadv/gopher-lua-libs/cmd"
	crypto "github.com/vadv/gopher-lua-libs/crypto"
	db "github.com/vadv/gopher-lua-libs/db"
	filepath "github.com/vadv/gopher-lua-libs/filepath"
	goos "github.com/vadv/gopher-lua-libs/goos"
	http "github.com/vadv/gopher-lua-libs/http"
	humanize "github.com/vadv/gopher-lua-libs/humanize"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	ioutil "github.com/vadv/gopher-lua-libs/ioutil"
	json "github.com/vadv/gopher-lua-libs/json"
	log "github.com/vadv/gopher-lua-libs/log"
	pb "github.com/vadv/gopher-lua-libs/pb"
	prometheus "github.com/vadv/gopher-lua-libs/prometheus/client"
	regexp "github.com/vadv/gopher-lua-libs/regexp"
	storage "github.com/vadv/gopher-lua-libs/storage"
	strings "github.com/vadv/gopher-lua-libs/strings"
	tac "github.com/vadv/gopher-lua-libs/tac"
	tcp "github.com/vadv/gopher-lua-libs/tcp"
	telegram "github.com/vadv/gopher-lua-libs/telegram"
	template "github.com/vadv/gopher-lua-libs/template"
	time "github.com/vadv/gopher-lua-libs/time"
	xmlpath "github.com/vadv/gopher-lua-libs/xmlpath"
	yaml "github.com/vadv/gopher-lua-libs/yaml"
	zabbix "github.com/vadv/gopher-lua-libs/zabbix"

	lua "github.com/yuin/gopher-lua"
)

type luaPlugin struct {
	sync.Mutex
	state      *lua.LState
	cancelFunc context.CancelFunc
	running    bool
	error      error
	body       *string
	filename   *string
	jobPayload *string
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
	// preload all
	filepath.Preload(state)
	http.Preload(state)
	inspect.Preload(state)
	ioutil.Preload(state)
	json.Preload(state)
	regexp.Preload(state)
	strings.Preload(state)
	tac.Preload(state)
	tcp.Preload(state)
	time.Preload(state)
	xmlpath.Preload(state)
	yaml.Preload(state)
	zabbix.Preload(state)
	telegram.Preload(state)
	storage.Preload(state)
	crypto.Preload(state)
	goos.Preload(state)
	humanize.Preload(state)
	db.Preload(state)
	chef.Preload(state)
	cmd.Preload(state)
	template.Preload(state)
	cloudwatch.Preload(state)
	log.Preload(state)
	prometheus.Preload(state)
	pb.Preload(state)
	stats.Preload(state)
	return state
}

func (p *luaPlugin) start() {
	p.Lock()
	state := NewPluginState()
	p.state = state
	p.error = nil
	p.running = true
	isBody := (p.filename == nil)
	if !(p.jobPayload == nil) {
		payload := *p.jobPayload
		state.SetGlobal(`payload`, lua.LString(payload))
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	p.cancelFunc = cancelFunc
	p.state.SetContext(ctx)
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

// DoString lua plugin.do_string(body) returns plugin_ud
func DoString(L *lua.LState) int {
	body := L.CheckString(1)
	p := &luaPlugin{body: &body}
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable(`plugin_ud`))
	L.Push(ud)
	return 1
}

// DoFile lua plugin.do_file(filename) returns plugin_ud
func DoFile(L *lua.LState) int {
	filename := L.CheckString(1)
	p := &luaPlugin{filename: &filename}
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
	p := &luaPlugin{filename: &filename, jobPayload: &payload}
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
	p := &luaPlugin{body: &body, jobPayload: &payload}
	ud := L.NewUserData()
	ud.Value = p
	L.SetMetatable(ud, L.GetTypeMetatable(`plugin_ud`))
	L.Push(ud)
	return 1
}

// Run lua plugin_ud:run()
func Run(L *lua.LState) int {
	p := checkPlugin(L, 1)
	go p.start()
	return 0
}

// IsRunning lua plugin_ud:is_running()
func IsRunning(L *lua.LState) int {
	p := checkPlugin(L, 1)
	L.Push(lua.LBool(p.getRunning()))
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
