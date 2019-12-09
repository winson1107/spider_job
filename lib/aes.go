package lib

import (
	"crypto/cipher"
	"crypto/aes"
	"log"
	"encoding/base64"
)

const (
	aesKey = "abcdefgh12345678abcdefgh"
	aesIv = "abcdabcd12345678"
)


//加密
func AesEny(plaintext []byte) string {
	var(
		block cipher.Block
		err error
	)
	//创建aes
	if block, err = aes.NewCipher([]byte(aesKey)); err != nil{
		log.Fatal(err)
	}
	//创建ctr
	stream := cipher.NewCTR(block, []byte(aesIv))
	//加密, src,dst 可以为同一个内存地址
	stream.XORKeyStream(plaintext, plaintext)
	return base64.StdEncoding.EncodeToString(plaintext)
}


//解密
func AesDec(ciptext []byte) (string,error) {
	var(
		block cipher.Block
		err error
	)
	ciptext,err = base64.StdEncoding.DecodeString(string(ciptext))
	if err != nil {
		return "", err
	}
	//创建aes
	if block, err = aes.NewCipher([]byte(aesKey)); err != nil{
		return "", err
	}
	//创建ctr
	stream := cipher.NewCTR(block,[]byte(aesIv))
	stream.XORKeyStream(ciptext, ciptext)
	return string(ciptext),nil
}
