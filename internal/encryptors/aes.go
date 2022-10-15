// Package encryptors includes all custom encryption-related functionality.
package encryptors

import (
	"crypto/aes"
	"encoding/hex"
)

// key is the encryption key
var key = []byte("auth_cipher_key0")

// AESEncrypt transforms original string into the encrypted one using the AES algorithm.
func AESEncrypt(data string) (string, error) {
	src := []byte(data[:aes.BlockSize])
	enc, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	dst := make([]byte, aes.BlockSize)
	enc.Encrypt(dst, src)
	return hex.EncodeToString(dst), nil
}

// AESDecrypt transforms the AES-encrypted string into the slice of bytes, containing the original value.
func AESDecrypt(data string) ([]byte, error) {
	src, err := hex.DecodeString(data)
	if err != nil {
		return nil, err
	}

	dec, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	dst := make([]byte, aes.BlockSize)
	dec.Decrypt(dst, src)
	return dst, nil
}
