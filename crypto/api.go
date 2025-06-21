// Package crypto implements golang package crypto functionality for lua.
package crypto

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// MD5 lua crypto.md5(string) return string
func MD5(L *lua.LState) int {
	str := L.CheckString(1)
	hash := md5.Sum([]byte(str))
	L.Push(lua.LString(fmt.Sprintf("%x", hash)))
	return 1
}

// SHA256 lua crypto.sha256(string) return string
func SHA256(L *lua.LState) int {
	str := L.CheckString(1)
	hash := sha256.Sum256([]byte(str))
	L.Push(lua.LString(fmt.Sprintf("%x", hash)))
	return 1
}

// AESEncrypt implements AES encryption in Lua.
func AESEncrypt(l *lua.LState) int {
	m, key, iv, data, err := decodeParams(l)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(fmt.Sprintf("failed to decode params: %v", err)))
		return 2
	}

	enc, err := encryptAES(mode(m), key, iv, data)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(fmt.Sprintf("failed to encrypt: %v", err)))
		return 2
	}
	l.Push(lua.LString(hex.EncodeToString(enc)))
	return 1
}

// AESDecrypt implement AES decryption in Lua.
func AESDecrypt(l *lua.LState) int {
	m, key, iv, data, err := decodeParams(l)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(fmt.Sprintf("failed to decode params: %v", err)))
		return 2
	}

	dec, err := decryptAES(mode(m), key, iv, data)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(fmt.Sprintf("failed to decrypt: %v", err)))
		return 2
	}

	l.Push(lua.LString(hex.EncodeToString(dec)))
	return 1
}
