// Package crypto implements golang package crypto functionality for lua.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
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
	key, iv, data, err := decodeParams(l)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(fmt.Sprintf("failed to decode params: %v", err)))
		return 2
	}

	enc, err := encrypt(key, iv, data)
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
	key, iv, data, err := decodeParams(l)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(fmt.Sprintf("failed to decode params: %v", err)))
		return 2
	}

	dec, err := decrypt(key, iv, data)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(fmt.Sprintf("failed to decrypt: %v", err)))
		return 2
	}

	l.Push(lua.LString(hex.EncodeToString(dec)))
	return 1
}

func decodeParams(l *lua.LState) (key, iv, data []byte, err error) {
	keyStr := l.ToString(1)
	ivStr := l.ToString(2)
	dataStr := l.ToString(3)

	key, err = hex.DecodeString(keyStr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode key: %v", err)
	}

	iv, err = hex.DecodeString(ivStr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode IV: %v", err)
	}

	data, err = hex.DecodeString(dataStr)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to decode data: %v", err)
	}

	return key, iv, data, nil
}

// encrypt implements AES CBC-128 ENCRYPTION which requires 3 data fields
// 1. Key (16 bytes)
// 2. Initialization Vector (IV) (16 bytes)
// 3. Encrypted Data (16 bytes or length multiple a of 16)
// The encrypted data is divided into blocks of 16 bytes (128 bits) which then operated on with the IV and Key.
func encrypt(key []byte, iv []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(data)%blockSize != 0 {
		return nil, fmt.Errorf("payload length %d is not a multiple of AES block size %d", len(data), blockSize)
	}

	if len(iv) != blockSize {
		return nil, fmt.Errorf("size of the IV %d is not the same as block size %d", len(iv), blockSize)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	mode.CryptBlocks(encrypted, data)

	return encrypted, nil
}

// decrypt implements AES CBC-128 DECRYPTION which requires 3 data fields
// 1. Key (16 bytes)
// 2. Initialization Vector (IV) (16 bytes)
// 3. Encrypted Data (16 bytes or length multiple a of 16)
// The encrypted data is divided into blocks of 16 bytes (128 bits) which then operated on with the IV and Key.
func decrypt(key []byte, iv []byte, encrypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(encrypted)%blockSize != 0 {
		return nil, fmt.Errorf("encrypted payload length %d is not a multiple of AES block size %d", len(encrypted), blockSize)
	}

	if len(iv) != blockSize {
		return nil, fmt.Errorf("size of the IV %d is not the same as block size %d", len(iv), blockSize)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	return decrypted, nil
}
