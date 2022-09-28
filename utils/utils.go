package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
)

func HandleError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func Hash(i interface{}) string {
	s := fmt.Sprint(i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}
