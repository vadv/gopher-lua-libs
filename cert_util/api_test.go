package cert_util

import (
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

	state := lua.NewState()
	Preload(state)
	if err := state.DoFile("./test/test_api.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
