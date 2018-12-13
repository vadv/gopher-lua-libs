package tcp

import (
	"net"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func runPingPongServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		handleTCPClient(conn)
	}
}

func handleTCPClient(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			return
		}
		data := buf[0:count]
		if string(data) == "ping" {
			conn.Write([]byte("pong\n"))
		} else {
			conn.Write([]byte("unknown\n"))
		}
	}
}

func TestApi(t *testing.T) {

	go runPingPongServer(":12345")
	time.Sleep(time.Second)

	state := lua.NewState()
	Preload(state)
	if err := state.DoFile("./test/test_api.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
