package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"reflect"
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

func GetStructValues(i interface{}) []interface{} {
	rv := reflect.ValueOf(i)
	args := make([]interface{}, 0)
	for i := 0; i < rv.NumField(); i++ {
		args = append(args, rv.Field(i).Interface())
	}
	return args
}
