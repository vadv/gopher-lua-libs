// Package cert_utils provides utility for checking ssl-certs in lua.
package cert_util

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// NotAfter lua cert.not_after(hostname, <ip>) returns (unixts cert_not_after, err)
func NotAfter(L *lua.LState) int {
	serverName, address := L.CheckString(1), ""
	if L.GetTop() > 1 {
		address = L.CheckString(2)
	} else {
		address = fmt.Sprintf("%s:443", serverName)
	}
	conn, err := net.DialTimeout(`tcp`, address, 5*time.Second)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	client := tls.Client(conn, &tls.Config{ServerName: serverName, InsecureSkipVerify: true})
	handshakeErr := client.Handshake()

	var minNotAfter *time.Time
	for _, cert := range client.ConnectionState().PeerCertificates {
		if minNotAfter == nil || minNotAfter.Unix() > cert.NotAfter.Unix() {
			minNotAfter = &cert.NotAfter
		}
	}

	if minNotAfter == nil {
		L.Push(lua.LNil)
		L.Push(lua.LString("certs not found"))
		return 2
	}

	L.Push(lua.LNumber(minNotAfter.Unix()))
	if handshakeErr == nil {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(handshakeErr.Error()))
	}
	return 2
}
