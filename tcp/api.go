// Package tcp implements tcp client lib for lua.
package tcp

import (
	lio "github.com/vadv/gopher-lua-libs/io"
	"net"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const (
	// timeout for dial
	DefaultDialTimeout = 5 * time.Second
	// timeout for write
	DefaultWriteTimeout = time.Second
	// timeout for read
	DefaultReadTimeout = time.Second
	// timeout for close
	DefaultCloseTimeout = time.Second
)

type luaTCPClient struct {
	net.Conn
	address string
}

func (c *luaTCPClient) connect() error {
	conn, err := net.DialTimeout("tcp", c.address, DefaultDialTimeout)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func checkLuaTCPClient(L *lua.LState, n int) *luaTCPClient {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*luaTCPClient); ok {
		return v
	}
	L.ArgError(n, "tcp connection expected")
	return nil
}

// Open lua tcp.open(string) returns (tcp_client_ud, err)
func Open(L *lua.LState) int {
	addr := L.CheckString(1)
	t := &luaTCPClient{address: addr}
	if err := t.connect(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = t
	L.SetMetatable(ud, L.GetTypeMetatable("tcp_client_ud"))
	L.Push(ud)
	return 1
}

// Write lua tcp_client_ud:write() returns err
func Write(L *lua.LState) int {
	conn := checkLuaTCPClient(L, 1)
	_ = conn.SetWriteDeadline(time.Now().Add(DefaultWriteTimeout))
	return lio.IOWriterWrite(L)
}

// Read lua tcp_client_ud:read(max_size_int) returns (string, err)
func Read(L *lua.LState) int {
	conn := checkLuaTCPClient(L, 1)
	// Backward compatibility for callers that don't pass a length
	if L.GetTop() < 2 {
		L.Push(lua.LNumber(1024))
	}
	_ = conn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	return lio.IOReaderRead(L)
}

// Close lua tcp_client_ud:close()
func Close(L *lua.LState) int {
	conn := checkLuaTCPClient(L, 1)
	_ = conn.SetDeadline(time.Now().Add(DefaultCloseTimeout))
	return lio.IOWriterClose(L)
}
