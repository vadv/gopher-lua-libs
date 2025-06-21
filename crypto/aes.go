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
		ciphertext := aesGCM.Seal(nil, init, plaintext, nil)
		return ciphertext, nil
	case CBC:
		padded := pad(plaintext, aes.BlockSize)
		mode := cipher.NewCBCEncrypter(block, init)
		ciphertext := make([]byte, len(padded))
		mode.CryptBlocks(ciphertext, padded)
		return ciphertext, nil
	case CTR:
		stream := cipher.NewCTR(block, init)
		ciphertext := make([]byte, len(plaintext))
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	default:
		return nil, fmt.Errorf("unsupported mode: %d", m)
	}
}

// decryptAES implements AES decryption given mode, key, cyphertext, and init value.
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

// encrypt implements AES CBC-128 ENCRYPTION which requires 3 data fields
// // 1. Key (16 bytes)
// // 2. Initialization Vector (IV) (16 bytes)
// // 3. Encrypted Data (16 bytes or length multiple a of 16)
// // The encrypted data is divided into blocks of 16 bytes (128 bits) which then operated on with the IV and Key.
// func encrypt(key []byte, iv []byte, data []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	blockSize := block.BlockSize()
// 	if len(data)%blockSize != 0 {
// 		return nil, fmt.Errorf("payload length %d is not a multiple of AES block size %d", len(data), blockSize)
// 	}

// 	if len(iv) != blockSize {
// 		return nil, fmt.Errorf("size of the IV %d is not the same as block size %d", len(iv), blockSize)
// 	}

// 	mode := cipher.NewCBCEncrypter(block, iv)
// 	encrypted := make([]byte, len(data))
// 	mode.CryptBlocks(encrypted, data)

// 	return encrypted, nil
// }

// // decrypt implements AES CBC-128 DECRYPTION which requires 3 data fields
// // 1. Key (16 bytes)
// // 2. Initialization Vector (IV) (16 bytes)
// // 3. Encrypted Data (16 bytes or length multiple a of 16)
// // The encrypted data is divided into blocks of 16 bytes (128 bits) which then operated on with the IV and Key.
// func decrypt(key []byte, iv []byte, encrypted []byte) ([]byte, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}

// 	blockSize := block.BlockSize()
// 	if len(encrypted)%blockSize != 0 {
// 		return nil, fmt.Errorf("encrypted payload length %d is not a multiple of AES block size %d", len(encrypted), blockSize)
// 	}

// 	if len(iv) != blockSize {
// 		return nil, fmt.Errorf("size of the IV %d is not the same as block size %d", len(iv), blockSize)
// 	}

// 	mode := cipher.NewCBCDecrypter(block, iv)
// 	decrypted := make([]byte, len(encrypted))
// 	mode.CryptBlocks(decrypted, encrypted)

// 	return decrypted, nil
// }
