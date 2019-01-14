package http

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

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
func (s *luaServer) serve() {
	s.err = http.Serve(s.Listener, s)
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
	go server.serve()
	ud := L.NewUserData()
	ud.Value = server
	L.SetMetatable(ud, L.GetTypeMetatable("http_server_ud"))
	L.Push(ud)
	return 1
}

// NewRequest return lua table with http.Request representation
func NewRequest(L *lua.LState, req *http.Request) *lua.LTable {
	luaRequest := L.NewTable()
	luaRequest.RawSetString(`host`, lua.LString(req.Host))
	luaRequest.RawSetString(`method`, lua.LString(req.Method))
	luaRequest.RawSetString(`referer`, lua.LString(req.Referer()))
	luaRequest.RawSetString(`proto`, lua.LString(req.Proto))
	luaRequest.RawSetString(`user_agent`, lua.LString(req.UserAgent()))
	if req.URL != nil && len(req.URL.Query()) > 0 {
		query := L.NewTable()
		for k, v := range req.URL.Query() {
			if len(v) > 0 {
				query.RawSetString(k, lua.LString(v[0]))
			}
		}
		luaRequest.RawSetString(`query`, query)
	}
	if len(req.Header) > 0 {
		headers := L.NewTable()
		for k, v := range req.Header {
			if len(v) > 0 {
				headers.RawSetString(k, lua.LString(v[0]))
			}
		}
		luaRequest.RawSetString(`headers`, headers)
	}
	luaRequest.RawSetString(`path`, lua.LString(req.URL.Path))
	luaRequest.RawSetString(`raw_path`, lua.LString(req.URL.RawPath))
	luaRequest.RawSetString(`raw_query`, lua.LString(req.URL.RawQuery))
	luaRequest.RawSetString(`request_uri`, lua.LString(req.RequestURI))
	luaRequest.RawSetString(`remote_addr`, lua.LString(req.RemoteAddr))
	return luaRequest
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
	http_util.Preload(state)

	http_server_response_writer_ud := state.NewTypeMetatable(`http_server_response_writer_ud`)
	state.SetGlobal(`http_server_response_writer_ud`, http_server_response_writer_ud)
	state.SetField(http_server_response_writer_ud, "__index", state.SetFuncs(state.NewTable(), map[string]lua.LGFunction{
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
