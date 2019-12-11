package lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

const (
	aesKey = "1234567890123456"
	aesIv  = "Impassphrasegood"
)

func AESEncrypt(src string) (string, error) {
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		fmt.Println("key error1", err)
		return "", err
	}
	if src == "" {
		return "", fmt.Errorf("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, []byte(aesIv))
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AESDecrypt(crypt string) ([]byte, error) {
	dst, err := base64.StdEncoding.DecodeString(crypt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return nil, err
	}
	if len(dst) == 0 {
		return nil, fmt.Errorf("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, []byte(aesIv))
	decrypted := make([]byte, len(dst))
	ecb.CryptBlocks(decrypted, dst)
	return PKCS5Trimming(decrypted), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
