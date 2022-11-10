package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	lua_json "github.com/vadv/gopher-lua-libs/json"
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

type LuaClient struct {
	*http.Client
	userAgent       string
	basicAuthUser   *string
	basicAuthPasswd *string
	headers         map[string]string
	debug           bool
}

// newLuaClient() returns new LuaClient
func newLuaClient() *LuaClient {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	return &LuaClient{Client: &http.Client{Jar: jar}}
}

func (client *LuaClient) updateRequest(req *http.Request) {
	// set basic auth
	if client.basicAuthUser != nil && client.basicAuthPasswd != nil {
		username, password := client.basicAuthUser, client.basicAuthPasswd
		req.SetBasicAuth(*username, *password)
	}
	// set user agent
	req.Header.Set(`User-Agent`, client.userAgent)
	// set headers
	if client.headers != nil {
		for k, v := range client.headers {
			req.Header.Set(k, v)
		}
	}
}

// DoRequest() process request with needed settings for request
func (client *LuaClient) DoRequest(req *http.Request) (*http.Response, error) {
	client.updateRequest(req)
	if client.debug {
		dump, _ := httputil.DumpRequestOut(req, true)
		log.Printf("[DEBUG] send request:\n%s\n", dump)
	}
	return client.Do(req)
}

// PostFormRequest() process Form
func (client *LuaClient) PostFormRequest(url string, data url.Values) (*http.Response, error) {
	return client.PostForm(url, data)
}

func checkClient(L *lua.LState) *LuaClient {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*LuaClient); ok {
		return v
	}
	L.ArgError(1, "http client expected")
	return nil
}

// http.client(config) returns (user data, error)
// config table:
//
//	{
//	  proxy="http(s)://<user>:<password>@host:<port>",
//	  timeout= 10,
//	  insecure_ssl=false,
//	  user_agent = "gopher-lua",
//	  basic_auth_user = "",
//	  basic_auth_password = "",
//	  headers = {"key"="value"},
//	  debug = false,
//	}
func New(L *lua.LState) int {
	var config *lua.LTable
	if L.GetTop() > 0 {
		config = L.CheckTable(1)
	}
	client := &LuaClient{Client: &http.Client{Timeout: DefaultTimeout}, userAgent: DefaultUserAgent}
	tlsConfig := &tls.Config{}
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
		// Client Cert and Key go together and handling in loop is challenging - just pull them out here
		clientPublicCertPEMFile := L.GetField(config, `client_public_cert_pem_file`)
		clientPrivateKeyPemFile := L.GetField(config, `client_private_key_pem_file`)
		if clientPublicCertPEMFile != lua.LNil && clientPrivateKeyPemFile != lua.LNil {
			if _, ok := clientPublicCertPEMFile.(lua.LString); !ok {
				L.ArgError(1, "client_public_cert_pem_file must be string")
			}
			if _, ok := clientPrivateKeyPemFile.(lua.LString); !ok {
				L.ArgError(1, "client_private_key_pem_file must be string")
			}
			clientCert, err := tls.LoadX509KeyPair(clientPublicCertPEMFile.String(), clientPrivateKeyPemFile.String())
			if err != nil {
				L.RaiseError("error loading client certificate from %s and %s: %v",
					clientPublicCertPEMFile, clientPrivateKeyPemFile, err)
			}
			tlsConfig.Certificates = []tls.Certificate{clientCert}
			transport.TLSClientConfig = tlsConfig
		}
		config.ForEach(func(k lua.LValue, v lua.LValue) {
			switch k.String() {
			// parse timeout
			case `timeout`:
				if value, ok := v.(lua.LNumber); ok {
					client.Timeout = time.Duration(value) * time.Second
				} else {
					L.ArgError(1, "timeout must be number")
				}
			// parse proxy
			case `proxy`:
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
			// parse insecure_ssl
			case `insecure_ssl`:
				if value, ok := v.(lua.LBool); ok {
					tlsConfig.InsecureSkipVerify = bool(value)
					transport.TLSClientConfig = tlsConfig
				} else {
					L.ArgError(1, "insecure_ssl must be bool")
				}
			// parse root_cas
			case `root_cas_pem_file`:
				if value, ok := v.(lua.LString); ok {
					pemData, err := ioutil.ReadFile(string(value))
					if err != nil {
						L.RaiseError("error loading root_cas_pem_file from %s: %v", value, err)
					}
					tlsConfig.RootCAs = x509.NewCertPool()
					tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
					transport.TLSClientConfig = tlsConfig
				} else {
					L.ArgError(1, "root_cas_pem_file must be string")
				}
			// parse user_agent
			case `user_agent`:
				if _, ok := v.(lua.LString); ok {
					client.userAgent = v.String()
				} else {
					L.ArgError(1, "user_agent must be string")
				}
			// parse basic_auth_user
			case `basic_auth_user`:
				if _, ok := v.(lua.LString); ok {
					user := v.String()
					client.basicAuthUser = &user
				} else {
					L.ArgError(1, "basic_auth_user must be string")
				}
			// parse basic_auth_password
			case `basic_auth_password`:
				if _, ok := v.(lua.LString); ok {
					password := v.String()
					client.basicAuthPasswd = &password
				} else {
					L.ArgError(1, "basic_auth_password must be string")
				}
			// parse debug
			case `debug`:
				if value, ok := v.(lua.LBool); ok {
					client.debug = bool(value)
				} else {
					L.ArgError(1, "debug must be bool")
				}
			// parse headers
			case `headers`:
				if tbl, ok := v.(*lua.LTable); ok {
					headers := make(map[string]string, 0)
					data, err := lua_json.ValueEncode(tbl)
					if err != nil {
						L.ArgError(1, "headers must be table of key-values string")
					}
					if err := json.Unmarshal(data, &headers); err != nil {
						L.ArgError(1, "headers must be table of key-values string")
					}
					client.headers = headers
				} else {
					L.ArgError(1, "headers must be table")
				}
			}
		})
	}

	// cookie support
	jar, _ := cookiejar.New(&cookiejar.Options{})
	client.Jar = jar
	client.Transport = transport
	ud := L.NewUserData()
	ud.Value = client
	L.SetMetatable(ud, L.GetTypeMetatable("http_client_ud"))
	L.Push(ud)
	return 1
}
