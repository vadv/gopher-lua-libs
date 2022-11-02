package http_test

import (
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vadv/gopher-lua-libs/tests"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	lua_http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	plugin "github.com/vadv/gopher-lua-libs/plugin"
	lua_time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

func basicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}
		handler(w, r)
	}
}

func httpCheckHeaders(w http.ResponseWriter, r *http.Request) {
	if header := r.Header.Get(`simple_header`); header != `check_header` {
		w.WriteHeader(412)
		return
	}
	w.Write([]byte(`OK`))
}

func httpCheckUserAgent(w http.ResponseWriter, r *http.Request) {
	if ua := r.UserAgent(); ua != `check_ua` {
		w.WriteHeader(412)
		return
	}
	w.Write([]byte(`OK`))
}

func getFormFile(r *http.Request, key, filename string) (err error) {
	file, header, err := r.FormFile(key)
	if err != nil {
		return err
	}
	defer file.Close()
	if header.Filename != filename {
		return fmt.Errorf("bad filename, get: %s except: %s\n", header.Filename, filename)
	}
	return nil
}

func httpUploadFile(w http.ResponseWriter, r *http.Request) {
	err := getFormFile(r, "file", "test.txt")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	w.Write([]byte(`OK`))
}

func httpUploadFileWithFields(w http.ResponseWriter, r *http.Request) {
	err := getFormFile(r, "file", "test.txt")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if r.FormValue("foo") != "bar" {
		w.WriteHeader(400)
		return
	}
	w.Write([]byte(`OK`))
}

func httpUploadMultipleFile(w http.ResponseWriter, r *http.Request) {
	err := getFormFile(r, "file", "test.txt")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = getFormFile(r, "file1", "test1.txt")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	w.Write([]byte(`OK`))
}

func httpUploadMultipleFileWithFields(w http.ResponseWriter, r *http.Request) {
	err := getFormFile(r, "file", "test.txt")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	err = getFormFile(r, "file1", "test1.txt")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if r.FormValue("foo") != "bar" {
		w.WriteHeader(400)
		return
	}
	w.Write([]byte(`OK`))
}

func httpRouterGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`OK`))
}

func httpRouterGetTimeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	w.Write([]byte(`OK`))
}

func runHttp(addr string) {
	err := http.ListenAndServe(addr, nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func runHttps(addr string) {
	err := http.ListenAndServeTLS(addr, "./test/server.crt", "./test/server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func request(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if strings.HasPrefix(url, string(body)) {
		return fmt.Errorf("bad url, get: %s except: %s\n", body, url)
	}
	return nil
}

func manyRequest(addr string) *errgroup.Group {
	eg := &errgroup.Group{}
	eg.Go(func() error {
		time.Sleep(5 * time.Second)
		for count := 0; count < 10; count++ {
			url := fmt.Sprintf("%s/%s?d=%d", addr, "url", count)
			func(url string) {
				eg.Go(func() error {
					return request(url)
				})
			}(url)
		}
		return nil
	})
	return eg
}

func TestApi(t *testing.T) {

	http.HandleFunc("/get", httpRouterGet)
	http.HandleFunc("/getBasicAuth", basicAuth(httpRouterGet, "admin", "123456", "Please enter your username and password for this site"))
	http.HandleFunc("/timeout", httpRouterGetTimeout)
	http.HandleFunc("/checkHeader", httpCheckHeaders)
	http.HandleFunc("/checkUserAgent", httpCheckUserAgent)
	http.HandleFunc("/upload", httpUploadFile)
	http.HandleFunc("/uploadWithFields", httpUploadFileWithFields)
	http.HandleFunc("/uploadMultiple", httpUploadMultipleFile)
	http.HandleFunc("/uploadMultipleWithFields", httpUploadMultipleFileWithFields)

	go runHttp(":1111")
	go runHttps(":1112")
	time.Sleep(time.Second)

	state := lua.NewState()
	defer state.Close()

	lua_http.Preload(state)
	lua_time.Preload(state)
	inspect.Preload(state)
	plugin.Preload(state)

	t.Run("test_client", func(t *testing.T) {
		assert.NoError(t, state.DoFile("./test/test_client.lua"))
	})

	t.Run("test_server_accept", func(t *testing.T) {
		eg := manyRequest("http://127.0.0.1:1113")
		assert.NoError(t, state.DoFile("./test/test_server_accept.lua"))
		assert.NoError(t, eg.Wait())
	})

	t.Run("test_server_handle", func(t *testing.T) {
		eg := manyRequest("http://127.0.0.1:2113")
		assert.NoError(t, state.DoFile("./test/test_server_handle.lua"))
		assert.NoError(t, eg.Wait())
	})

	t.Run("test_server_accept_stop", func(t *testing.T) {
		assert.NoError(t, state.DoFile("./test/test_server_accept_stop.lua"))
	})

	t.Run("test_serve_static", func(t *testing.T) {
		assert.NoError(t, state.DoFile("./test/test_serve_static.lua"))
	})
}

func TestMTLS(t *testing.T) {
	s := httptest.NewUnstartedServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(writer, "OK\n")
	}))
	defer s.Close()
	serverCert, err := tls.LoadX509KeyPair("test/data/test.cert.pem", "test/data/test.key.pem")
	require.NoError(t, err)
	caData, err := ioutil.ReadFile("test/data/test.cert.pem")
	require.NoError(t, err)
	cas := x509.NewCertPool()
	cas.AppendCertsFromPEM(caData)
	s.TLS = &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    cas,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	s.StartTLS()

	preload := tests.SeveralPreloadFuncs(
		lua_http.Preload,
		lua_time.Preload,
		inspect.Preload,
		plugin.Preload,
		func(L *lua.LState) {
			// Attach the server URL to the testing object
			L.SetGlobal("tURL", lua.LString(s.URL))
		},
	)
	assert.NotZero(t, tests.RunLuaTestFile(t, preload, "test/test_api.lua"))
}
