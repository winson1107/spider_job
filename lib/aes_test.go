package lib

import (
	"testing"
	"log"
)

func TestAesDec(t *testing.T) {
	s,_ := AesDec([]byte("YlnC/RzgAxlWlAg="))
	log.Println(s)
}
func TestAesEny(t *testing.T) {
	s := AesEny([]byte("test"))
	log.Println(s)
}
