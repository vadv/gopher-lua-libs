package http_test

import (
	"crypto/subtle"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	client := &http.Client{}
	resp, err := client.Get(url)
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

func manyRequest(addr string) {
	time.Sleep(5 * time.Second)
	count := 0
	for {
		if count > 10 {
			break
		}
		url := fmt.Sprintf("%s/%s?d=%d", addr, "url", count)
		go func(url string) {
			if err := request(url); err != nil {
				panic(err)
			}
		}(url)
		count++
	}
}

func TestApi(t *testing.T) {

	http.HandleFunc("/get", httpRouterGet)
	http.HandleFunc("/getBasicAuth", basicAuth(httpRouterGet, "admin", "123456", "Please enter your username and password for this site"))
	http.HandleFunc("/timeout", httpRouterGetTimeout)
	http.HandleFunc("/checkHeader", httpCheckHeaders)
	http.HandleFunc("/checkUserAgent", httpCheckUserAgent)

	go runHttp(":1111")
	go runHttps(":1112")
	time.Sleep(time.Second)

	state := lua.NewState()
	lua_http.Preload(state)
	lua_time.Preload(state)
	inspect.Preload(state)
	plugin.Preload(state)

	if err := state.DoFile("./test/test_client.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}

	go manyRequest("http://127.0.0.1:1113")
	if err := state.DoFile("./test/test_server_accept.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}

	go manyRequest("http://127.0.0.1:2113")
	if err := state.DoFile("./test/test_server_handle.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
