package http

import (
	"crypto/subtle"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

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

func TestApi(t *testing.T) {

	http.HandleFunc("/get", httpRouterGet)
	http.HandleFunc("/getBasicAuth", basicAuth(httpRouterGet, "admin", "123456", "Please enter your username and password for this site"))
	http.HandleFunc("/timeout", httpRouterGetTimeout)
	http.HandleFunc("/checkHeader", httpCheckHeaders)
	http.HandleFunc("/checkUserAgent", httpCheckUserAgent)

	go runHttp(":1111")
	go runHttps(":1112")
	time.Sleep(time.Second)

	data, err := ioutil.ReadFile("./test/test_api.lua")
	if err != nil {
		t.Fatalf("%s\n", err.Error())
	}
	state := lua.NewState()
	Preload(state)
	if err := state.DoString(string(data)); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
