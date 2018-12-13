package cert_util

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func runHttps(addr string) {
	err := http.ListenAndServeTLS(addr, "./test/cert.pem", "./test/key.pem", nil)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func httpRouterGet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`OK`))
}

func TestApi(t *testing.T) {

	http.HandleFunc("/get", httpRouterGet)
	go runHttps(":1443")
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
