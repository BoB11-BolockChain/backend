package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/backend/auth"
	"github.com/backend/board"
	"github.com/backend/caldera"
	"github.com/backend/challenges"
	"github.com/backend/create"
	"github.com/backend/utils"
	"github.com/gorilla/mux"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		next.ServeHTTP(w, r)
	})
}

func optionMethodBanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		} else {
			w.Write([]byte("message from option block middleware"))
		}
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	d := make(map[string]interface{})

	json.NewEncoder(w).Encode(d)
}

func getabs(w http.ResponseWriter, r *http.Request) {
	addr := "http://www.pdxf.tk:8888/api/v2/abilities"
	req, err := http.NewRequest("GET", addr, nil)
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")

	client := http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	utils.HandleError(err)

	fmt.Fprint(w, string(b))
}

func Start(port int) {
	addr := fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.Use(corsMiddleware)
	router.Use(optionMethodBanMiddleware)

	router.HandleFunc("/", hello)
	router.HandleFunc("/abilities", getabs)

	router.HandleFunc("/signin", auth.SignIn)
	router.HandleFunc("/signup", auth.SignUp)
	router.HandleFunc("/logout", auth.Logout)
	router.HandleFunc("/welcome", auth.Welcome)
	router.HandleFunc("/profile", auth.UserInfo)

	router.HandleFunc("/challenges", challenges.ChInfo)
	router.HandleFunc("/info", challenges.ViewInfo)

	router.HandleFunc("/createchallenges", challenges.InsertData)
	// router.HandleFunc("/createch2", challenges.InsertData2)

	router.HandleFunc("/getch", challenges.PrintData)
	router.HandleFunc("/basic", challenges.LoadBasic)

	router.HandleFunc("/docker", create.DockerRun)
	router.HandleFunc("/vagrant", create.VagrantRun)

	router.HandleFunc("/notification", board.Notification)
	router.HandleFunc("/noticreate", board.NotiCreate)
	router.HandleFunc("/notiedit", board.NotiEdit)

	router.HandleFunc("/operation", caldera.GetOperationId)

	router.HandleFunc("/dashboard", challenges.SocketEndpoint)
	router.HandleFunc("/createoperation", challenges.CreateOperation).Methods("POST")
	// 메소드 지정할수 잇내요^^;

	log.Fatal(http.ListenAndServe(addr, router))
}

func main() {
	var port int
	fmt.Println("사용할 포트 입력 (수정 : 3000, 성현 : 8000) : ")
	fmt.Scanf("%d", &port)
	Start(port)
}
