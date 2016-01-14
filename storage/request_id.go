package storage

import (
//"crypto/aes"
//"crypto/cipher"
//"crypto/rand"
//"encoding/hex"
//"fmt"
//"io"
//"log"
)

const (
	aesKey = "268561882327517864681035"
)

func encryptRequestId(reqId string) (string, error) {
	//block, err := aes.NewCipher([]byte(aesKey))
	//if err != nil {
	//	return "", err
	//}
	//ciphertext := make([]byte, aes.BlockSize+len(reqId))
	//iv := ciphertext[:aes.BlockSize]
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	return "", err
	//}

	//cfb := cipher.NewCFBEncrypter(block, iv)
	//cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(reqId))

	//log.Printf(string(ciphertext))
	//return hex.EncodeToString(ciphertext), nil
	return reqId, nil
}

func decryptRequestId(reqId string) (string, error) {
	//block, err := aes.NewCipher([]byte(aesKey))
	//if err != nil {
	//	return "", err
	//}

	//byteReqId, err := hex.DecodeString(reqId)
	//if err != nil {
	//	return "", err
	//}
	//if len(byteReqId) < aes.BlockSize {
	//	return "", fmt.Errorf("Encrypted request id is too short")
	//}

	//iv := byteReqId[:aes.BlockSize]
	//ciphertext := byteReqId[aes.BlockSize:]

	//cfb := cipher.NewCFBDecrypter(block, iv)
	//cfb.XORKeyStream(ciphertext, ciphertext)
	//log.Printf(string(ciphertext))
	//return string(ciphertext), nil
	return reqId, nil
}
