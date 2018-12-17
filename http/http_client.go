package http

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"os"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const (
	// default http User Agent
	DefaultUserAgent = `gopher-lua`
	// default http timeout
	DefaultTimeout = 10 * time.Second
	// default don't ignore ssl
	insecureSkipVerify = false
)

type LuaHTTPClient struct {
	*http.Client
}

func checkClient(L *lua.LState) *LuaHTTPClient {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*LuaHTTPClient); ok {
		return v
	}
	L.ArgError(1, "http client excepted")
	return nil
}

// http.client(config) returns (user data, error)
// config table:
//   {
//     proxy="http(s)://<user>:<password>@host:<port>",
//     timeout= 10,
//     insecure_ssl=false,
//   }
func NewClient(L *lua.LState) int {
	var config *lua.LTable
	if L.GetTop() > 0 {
		config = L.CheckTable(1)
	}
	client := &LuaHTTPClient{Client: &http.Client{Timeout: DefaultTimeout}}
	transport := &http.Transport{}
	// parse env
	if proxyEnv := os.Getenv(`HTTP_PROXY`); proxyEnv != `` {
		proxyUrl, err := url.Parse(proxyEnv)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	transport.MaxIdleConns = 0
	transport.MaxIdleConnsPerHost = 1
	transport.IdleConnTimeout = DefaultTimeout
	// parse config
	if config != nil {
		config.ForEach(func(k lua.LValue, v lua.LValue) {
			// parse timeout
			if k.String() == `timeout` {
				if value, ok := v.(lua.LNumber); ok {
					client.Timeout = time.Duration(value) * time.Second
				} else {
					L.ArgError(1, "timeout must be number")
				}
			}
			// parse proxy
			if k.String() == `proxy` {
				if value, ok := v.(lua.LString); ok {
					proxyUrl, err := url.Parse(value.String())
					if err == nil {
						transport.Proxy = http.ProxyURL(proxyUrl)
					} else {
						L.ArgError(1, "http_proxy must be http(s)://<user>:<password>@host:<port>")
					}
				} else {
					L.ArgError(1, "http_proxy must be string")
				}
			}
			// parse insecure_ssl
			if k.String() == `insecure_ssl` {
				if value, ok := v.(lua.LBool); ok {
					transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: bool(value)}
				} else {
					L.ArgError(1, "insecure_ssl must be bool")
				}
			}
		})
	}

	client.Transport = transport
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("http_client_ud"))
	L.Push(ud)
	return 1
}
