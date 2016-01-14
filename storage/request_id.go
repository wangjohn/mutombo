package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const (
	aesKey = "268561882327517864681035"
)

func encryptRequestId(reqId string) (string, error) {
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(reqId))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(reqId))

	return string(ciphertext), nil
}

func decryptRequestId(reqId string) (string, error) {
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}

	byteReqId := []byte(reqId)
	if len(byteReqId) < aes.BlockSize {
		return "", fmt.Errorf("Encrypted request id is too short")
	}

	iv := byteReqId[:aes.BlockSize]
	ciphertext := byteReqId[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
