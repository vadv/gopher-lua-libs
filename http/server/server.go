package http

import (
	"crypto/tls"
	"crypto/x509"
	ioutil2 "io/ioutil"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/vadv/gopher-lua-libs/aws/cloudwatch"
	"github.com/vadv/gopher-lua-libs/cert_util"
	"github.com/vadv/gopher-lua-libs/chef"
	"github.com/vadv/gopher-lua-libs/cmd"
	"github.com/vadv/gopher-lua-libs/crypto"
	"github.com/vadv/gopher-lua-libs/db"
	"github.com/vadv/gopher-lua-libs/filepath"
	"github.com/vadv/gopher-lua-libs/goos"
	http_client "github.com/vadv/gopher-lua-libs/http/client"
	http_util "github.com/vadv/gopher-lua-libs/http/util"
	"github.com/vadv/gopher-lua-libs/humanize"
	"github.com/vadv/gopher-lua-libs/inspect"
	"github.com/vadv/gopher-lua-libs/ioutil"
	"github.com/vadv/gopher-lua-libs/json"
	lua_log "github.com/vadv/gopher-lua-libs/log"
	"github.com/vadv/gopher-lua-libs/prometheus/client"
	"github.com/vadv/gopher-lua-libs/regexp"
	"github.com/vadv/gopher-lua-libs/runtime"
	"github.com/vadv/gopher-lua-libs/storage"
	"github.com/vadv/gopher-lua-libs/strings"
	"github.com/vadv/gopher-lua-libs/tac"
	"github.com/vadv/gopher-lua-libs/tcp"
	"github.com/vadv/gopher-lua-libs/telegram"
	"github.com/vadv/gopher-lua-libs/template"
	lua_time "github.com/vadv/gopher-lua-libs/time"
	"github.com/vadv/gopher-lua-libs/xmlpath"
	"github.com/vadv/gopher-lua-libs/yaml"
	"github.com/vadv/gopher-lua-libs/zabbix"

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
	var tlsConfig *tls.Config
	bind := "127.0.0.1:0"
	switch bindOrTable := L.CheckAny(1).(type) {
	case lua.LString:
		bind = string(bindOrTable)
	case *lua.LTable:
		if addr, ok := L.GetField(bindOrTable, "addr").(lua.LString); ok {
			bind = string(addr)
		}
		serverPublicCertPEMFile := L.GetField(bindOrTable, `server_public_cert_pem_file`)
		serverPrivateKeyPemFile := L.GetField(bindOrTable, `server_private_key_pem_file`)
		if serverPublicCertPEMFile != lua.LNil && serverPrivateKeyPemFile != lua.LNil {
			serverCert, err := tls.LoadX509KeyPair(serverPublicCertPEMFile.String(), serverPrivateKeyPemFile.String())
			if err != nil {
				L.RaiseError("error loading server cert: %v", err)
			}
			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{serverCert},
			}

			clientAuth := L.GetField(bindOrTable, "client_auth")
			if clientAuth != lua.LNil {
				if _, ok := clientAuth.(lua.LString); !ok {
					L.ArgError(1, "client_auth should be a string")
				}
				switch clientAuth.String() {
				case "NoClientCert":
					tlsConfig.ClientAuth = tls.NoClientCert
				case "RequestClientCert":
					tlsConfig.ClientAuth = tls.RequestClientCert
				case "RequireAnyClientCert":
					tlsConfig.ClientAuth = tls.RequireAnyClientCert
				case "VerifyClientCertIfGiven":
					tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
				case "RequireAndVerifyClientCert":
					tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
				}
			}

			clientCAs := L.GetField(bindOrTable, "client_cas_pem_file")
			if clientCAs != lua.LNil {
				if _, ok := clientCAs.(lua.LString); !ok {
					L.ArgError(1, "client_cas_pem_file must be a string")
				}
				data, err := ioutil2.ReadFile(clientCAs.String())
				if err != nil {
					L.RaiseError("error reading %s: %v", clientCAs, err)
				}
				tlsConfig.ClientCAs = x509.NewCertPool()
				if !tlsConfig.ClientCAs.AppendCertsFromPEM(data) {
					L.RaiseError("no certs loaded from %s", clientCAs)
				}
			}
		}
	}
	l, err := net.Listen(`tcp`, bind)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	if tlsConfig != nil {
		l = tls.NewListener(l, tlsConfig)
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

// Addr returns the address if, for instance, one listens on :0
func Addr(L *lua.LState) int {
	s := checkServer(L, 1)
	L.Push(lua.LString(s.Listener.Addr().String()))
	return 1
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
				defer state.Close()
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

// HandleString lua http_server_ud:handle_string(body)
func HandleString(L *lua.LState) int {
	s := checkServer(L, 1)
	body := L.CheckString(2)
	for {
		select {
		case data := <-s.serveData:
			go func(sData *serveData, content string) {
				state := newHandlerState(sData)
				defer state.Close()
				if err := state.DoString(content); err != nil {
					log.Printf("[ERROR] handle: %s\n", err.Error())
					data.done <- true
					log.Printf("[ERROR] closed connection\n")
				}
			}(data, body)
		}
	}
	return 0
}

// HandleFunction lua http_server_ud:handle_function(func(response, request))
func HandleFunction(L *lua.LState) int {
	s := checkServer(L, 1)
	f := L.CheckFunction(2)
	if len(f.Upvalues) > 0 {
		L.ArgError(2, "cannot pass closures")
	}

	// Stash any args to pass to the function beyond response and request
	var args []lua.LValue
	top := L.GetTop()
	for i := 3; i <= top; i++ {
		args = append(args, L.Get(i))
	}

	for {
		select {
		case data := <-s.serveData:
			go func(sData *serveData) {
				state := newHandlerState(sData)
				defer state.Close()
				response := state.GetGlobal("response")
				request := state.GetGlobal("request")
				f := state.NewFunctionFromProto(f.Proto)
				state.Push(f)
				state.Push(response)
				state.Push(request)
				// Push any extra args
				for _, arg := range args {
					state.Push(arg)
				}
				if err := state.PCall(2+len(args), 0, nil); err != nil {
					log.Printf("[ERROR] handle: %s\n", err.Error())
					data.done <- true
					log.Printf("[ERROR] closed connection\n")
				}
				state.Pop(state.GetTop())
			}(data)
		}
	}
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
