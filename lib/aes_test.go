package lib

import (
	"testing"
	"fmt"
)
func TestAesEny(t *testing.T) {
	var plainText = "admin"

	encryptedData,_ := AESEncrypt(plainText)
	fmt.Println(encryptedData)
	decryptedText,_ := AESDecrypt(encryptedData)
	fmt.Println(string(decryptedText))
}