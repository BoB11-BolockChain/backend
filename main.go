package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/backend/auth"
	"github.com/backend/utils"
	"github.com/gorilla/mux"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func getabs(w http.ResponseWriter, r *http.Request) {
	addr := "http://domain:8888/api/v2/abilities"
	req, err := http.NewRequest("GET", addr, nil)
	utils.HandleError(err)

	req.Header.Add("admin", "admin123")

	client := http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	utils.HandleError(err)
	j, err := json.MarshalIndent(b, "", "  ")
	utils.HandleError(err)

	fmt.Fprint(w, string(j))
}

func Start(port int) {
	addr := fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)

	router.HandleFunc("/", hello)
	router.HandleFunc("/abilities", getabs)

	router.HandleFunc("/signin", auth.SignIn)

	log.Fatal(http.ListenAndServe(addr, router))
}

func main() {
	Start(8000)
}