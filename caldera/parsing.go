package caldera

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	// "github.com/utils"
)

type Operation struct {
	// jitter string
	// visibility int
	// auto_close bool
	// id string
	Objective interface{}
	// adversary []interface{}
	// autonomous int
	// group string
	// etc []interface{}

}

// type newId struct {
// 	Id string
// }

func GetOperationId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	// var op []Operation
	// var op Id
	// var op string
	url := "http://211.197.16.122:8888/api/v2/operations?include=id"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("KEY", "ADMIN123")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bytes, _ := io.ReadAll(resp.Body)
	str := string(bytes)
	// fmt.Println(resp.Body)
	// var jsonId newId
	// var val string
	val := (strings.Trim(str, "[]"))
	fmt.Println(val)
	// fmt.Println(jsonId)

	// slice := strings.Split(str, " ")

	// for _, str := range slice {
	// 	fmt.Println(str)
	// }

}
