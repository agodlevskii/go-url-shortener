package encryptors

import (
	"crypto/aes"
	"encoding/hex"
)

var key = []byte("auth_cipher_key0")

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
