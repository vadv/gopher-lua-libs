package http

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	cloudwatch "github.com/vadv/gopher-lua-libs/aws/cloudwatch"
	cert_util "github.com/vadv/gopher-lua-libs/cert_util"
	chef "github.com/vadv/gopher-lua-libs/chef"
	cmd "github.com/vadv/gopher-lua-libs/cmd"
	crypto "github.com/vadv/gopher-lua-libs/crypto"
	db "github.com/vadv/gopher-lua-libs/db"
	filepath "github.com/vadv/gopher-lua-libs/filepath"
	goos "github.com/vadv/gopher-lua-libs/goos"
	http_client "github.com/vadv/gopher-lua-libs/http/client"
	http_util "github.com/vadv/gopher-lua-libs/http/util"
	humanize "github.com/vadv/gopher-lua-libs/humanize"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	ioutil "github.com/vadv/gopher-lua-libs/ioutil"
	json "github.com/vadv/gopher-lua-libs/json"
	lua_log "github.com/vadv/gopher-lua-libs/log"
	"github.com/vadv/gopher-lua-libs/prometheus/client"
	regexp "github.com/vadv/gopher-lua-libs/regexp"
	runtime "github.com/vadv/gopher-lua-libs/runtime"
	storage "github.com/vadv/gopher-lua-libs/storage"
	strings "github.com/vadv/gopher-lua-libs/strings"
	tac "github.com/vadv/gopher-lua-libs/tac"
	tcp "github.com/vadv/gopher-lua-libs/tcp"
	telegram "github.com/vadv/gopher-lua-libs/telegram"
	template "github.com/vadv/gopher-lua-libs/template"
	lua_time "github.com/vadv/gopher-lua-libs/time"
	xmlpath "github.com/vadv/gopher-lua-libs/xmlpath"
	yaml "github.com/vadv/gopher-lua-libs/yaml"
	zabbix "github.com/vadv/gopher-lua-libs/zabbix"

	lua "github.com/yuin/gopher-lua"
)

type luaServer struct {
	*http.Server
	net.Listener
	sync.Mutex
	serveData chan *serveData
	err       error
}

type serveData struct {
	w    http.ResponseWriter
	req  *http.Request
	done chan bool
}

func checkServer(L *lua.LState, n int) *luaServer {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaServer); ok {
		return v
	}
	L.ArgError(n, "http server excepted")
	return nil
}

// run serve
func (s *luaServer) serve(L *lua.LState) {
	// start serve
	go func() {
		s.err = http.Serve(s.Listener, s)
	}()
	// process shutdown
	go func(s *luaServer) {
		ctx := L.Context()
		if ctx != nil {
			select {
			case <-ctx.Done():
				s.Listener.Close()
			}
		}
	}(s)
}

// http.server(bind, handler) returns (user data, error)
func New(L *lua.LState) int {
	bind := L.CheckAny(1).String()
	l, err := net.Listen(`tcp`, bind)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	server := &luaServer{
		Listener:  l,
		serveData: make(chan *serveData, 1),
	}
	server.serve(L)
	ud := L.NewUserData()
	ud.Value = server
	L.SetMetatable(ud, L.GetTypeMetatable("http_server_ud"))
	L.Push(ud)
	return 1
}

// Accept lua http_server_ud:accept() returns request_table, http_server_response_writer_ud
func Accept(L *lua.LState) int {
	s := checkServer(L, 1)
	select {
	case data := <-s.serveData:
		L.Push(NewRequest(L, data.req))
		L.Push(NewWriter(L, data.w, data.req, data.done))
		return 2
	}
}

func newHandlerState(data *serveData) *lua.LState {
	state := lua.NewState()

	lua_time.Preload(state)
	strings.Preload(state)
	filepath.Preload(state)
	ioutil.Preload(state)
	regexp.Preload(state)
	tac.Preload(state)
	inspect.Preload(state)
	yaml.Preload(state)
	cmd.Preload(state)
	json.Preload(state)
	tcp.Preload(state)
	xmlpath.Preload(state)
	db.Preload(state)
	cert_util.Preload(state)
	runtime.Preload(state)
	telegram.Preload(state)
	zabbix.Preload(state)
	crypto.Preload(state)
	goos.Preload(state)
	storage.Preload(state)
	humanize.Preload(state)
	chef.Preload(state)
	template.Preload(state)
	http_client.Preload(state)
	lua_log.Preload(state)
	cloudwatch.Preload(state)
	http_util.Preload(state)
	prometheus_client.Preload(state)

	httpServerResponseWriterUD := state.NewTypeMetatable(`http_server_response_writer_ud`)
	state.SetGlobal(`http_server_response_writer_ud`, httpServerResponseWriterUD)
	state.SetField(httpServerResponseWriterUD, "__index", state.SetFuncs(state.NewTable(), map[string]lua.LGFunction{
		"code":     HeaderCode,
		"header":   Header,
		"write":    Write,
		"redirect": Redirect,
		"done":     Done,
	}))

	state.SetGlobal("request", NewRequest(state, data.req))
	state.SetGlobal("response", NewWriter(state, data.w, data.req, data.done))

	return state
}

// HandleFile lua http_server_ud:handler_file(filename)
func HandleFile(L *lua.LState) int {
	s := checkServer(L, 1)
	file := L.CheckString(2)
	for {
		select {
		case data := <-s.serveData:
			go func(sData *serveData, filename string) {
				state := newHandlerState(data)
				if err := state.DoFile(filename); err != nil {
					log.Printf("[ERROR] handle file %s: %s\n", filename, err.Error())
					data.done <- true
					log.Printf("[ERROR] closed connection\n")
				}
			}(data, file)

		}
	}
	return 0
}

// HandleString lua http_server_ud:handler_string(body)
func HandleString(L *lua.LState) int {
	s := checkServer(L, 1)
	body := L.CheckString(2)
	select {
	case data := <-s.serveData:
		go func(sData *serveData, content string) {
			state := newHandlerState(sData)
			if err := state.DoString(content); err != nil {
				log.Printf("[ERROR] handle: %s\n", err.Error())
				data.done <- true
				log.Printf("[ERROR] closed connection\n")
			}
		}(data, body)
	}
	return 0
}

// ServeHTTP interface realisation
func (s *luaServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	doneChan := make(chan bool)
	data := &serveData{w: w, req: req, done: doneChan}
	// send data for lua
	s.serveData <- data

	// wait response from lua
	select {
	case <-doneChan:
		return
	case <-time.After(time.Minute):
		doneChan <- true
	}

}
