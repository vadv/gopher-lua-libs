// Package tcp implements tcp client lib for lua.
package tcp

import (
	"fmt"
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

// Open(): lua tcp.open(string) returns (tcp_client_ud, err)
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

// Write(): lua tcp_client_ud:write() returns err
func Write(L *lua.LState) int {
	conn := checkLuaTCPClient(L, 1)
	data := L.CheckString(2)
	conn.SetWriteDeadline(time.Now().Add(DefaultWriteTimeout))
	count, err := conn.Write([]byte(data))
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("write to `%s`: %s", conn.address, err.Error())))
		return 1
	}
	if count != len(data) {
		L.Push(lua.LString(fmt.Sprintf("write to `%s` get: %d except: %d", conn.address, count, len(data))))
		return 1
	}
	return 0
}

// Read(): lua tcp_client_ud:read(max_size_int) returns (string, err)
func Read(L *lua.LState) int {
	conn := checkLuaTCPClient(L, 1)
	count := int(1024)
	if L.GetTop() > 1 {
		count = int(L.CheckInt64(2))
		if count < 1 {
			L.ArgError(2, "must be > 1")
		}
	}
	buf := make([]byte, count)
	conn.SetReadDeadline(time.Now().Add(DefaultReadTimeout))
	count, err := conn.Read(buf)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(fmt.Sprintf("read from `%s`: %s", conn.address, err.Error())))
		return 2
	}
	line := string(buf[0:count])
	L.Push(lua.LString(line))
	return 1
}

// Close(): lua tcp_client_ud:close()
func Close(L *lua.LState) int {
	conn := checkLuaTCPClient(L, 1)
	conn.SetDeadline(time.Now().Add(DefaultCloseTimeout))
	if conn != nil {
		conn.Close()
	}
	return 0
}
