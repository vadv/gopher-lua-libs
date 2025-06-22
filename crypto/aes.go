package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type mode uint

const (
	GCM mode = iota + 1
	CBC
	CTR
)

func (m mode) String() string {
	switch m {
	case GCM:
		return "GCM"
	case CBC:
		return "CBC"
	case CTR:
		return "CTR"
	default:
		return "unknown"
	}
}

func decodeParams(l *lua.LState) (mode int, key, iv, data []byte, err error) {
	mode = l.ToInt(1)
	keyStr := l.ToString(2)
	ivStr := l.ToString(3)
	dataStr := l.ToString(4)

	key, err = hex.DecodeString(keyStr)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("failed to decode key: %v", err)
	}

	iv, err = hex.DecodeString(ivStr)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("failed to decode IV: %v", err)
	}

	data, err = hex.DecodeString(dataStr)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("failed to decode data: %v", err)
	}
	return mode, key, iv, data, nil
}

// encryptAES implements AES encryption given mode, key, plaintext, and init value.
// Init value is either initialization vector or nonce, depending on the mode.
func encryptAES(m mode, key, init, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	switch m {
	case GCM:
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		if len(init) != aesGCM.NonceSize() {
			return nil, fmt.Errorf("incorrect GCM nonce size: %d, expected: %d", len(init), aesGCM.NonceSize())
		}
		ciphertext := aesGCM.Seal(nil, init, plaintext, nil)
		return ciphertext, nil
	case CBC:
		if len(init) != block.BlockSize() {
			return nil, fmt.Errorf("invalid IV size: %d, expected: %d", len(init), block.BlockSize())
		}
		padded := pad(plaintext, aes.BlockSize)
		mode := cipher.NewCBCEncrypter(block, init)
		ciphertext := make([]byte, len(padded))
		mode.CryptBlocks(ciphertext, padded)
		return ciphertext, nil
	case CTR:
		l := len(init)
		if l != block.BlockSize() {
			return nil, fmt.Errorf("invalid IV size: %d, expected: %d", l, block.BlockSize())
		}
		stream := cipher.NewCTR(block, init)
		ciphertext := make([]byte, len(plaintext))
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	default:
		return nil, fmt.Errorf("unsupported mode: %d", m)
	}
}

// decryptAES implements AES decryption given mode, key, ciphertext, and init value.
// Init value is either initialization vector or nonce, depending on the mode.
func decryptAES(m mode, key, init, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	switch m {
	case GCM:
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		l := len(init)
		if l != aesGCM.NonceSize() {
			return nil, fmt.Errorf("incorrect GCM nonce size: %d, expected: %d", len(init), aesGCM.NonceSize())
		}
		plaintext, err := aesGCM.Open(nil, init, ciphertext, nil)
		if err != nil {
			return nil, err
		}
		return plaintext, nil
	case CBC:
		if len(ciphertext)%aes.BlockSize != 0 {
			return nil, fmt.Errorf("ciphertext is not a multiple of block size")
		}
		mode := cipher.NewCBCDecrypter(block, init)
		plaintext := make([]byte, len(ciphertext))
		mode.CryptBlocks(plaintext, ciphertext)
		return unpad(plaintext, aes.BlockSize)
	case CTR:
		stream := cipher.NewCTR(block, init)
		plaintext := make([]byte, len(ciphertext))
		stream.XORKeyStream(plaintext, ciphertext)
		return plaintext, nil
	default:
		return nil, fmt.Errorf("unsupported mode: %s", m)
	}
}

func pad(data []byte, blockSize int) []byte {
	padLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...)
}

func unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	padLen := int(data[len(data)-1])
	if padLen == 0 || padLen > blockSize {
		return nil, fmt.Errorf("invalid padding")
	}
	for i := 0; i < padLen; i++ {
		if data[len(data)-1-i] != byte(padLen) {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}
	return data[:len(data)-padLen], nil
}
