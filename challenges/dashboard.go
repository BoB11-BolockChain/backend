package challenges

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/backend/utils"
	"github.com/gorilla/websocket"
)

// caldera data. should be replaced with challenge db
const URL = "http://pdxf.tk:8888/api/v2/"
const ID = "9f4bd985-13b8-418b-be6c-5c2f5ba74829"

type OperationReport struct {
	Name       string `json:"name"`
	Start      string
	Host_group []interface{} `json:"host_group"`
	Steps      interface{}   `json:"steps"`
	Finish     bool
	Planner    string
	Adversary  interface{}
	Jitter     string
	Objectives interface{}
	Facts      []interface{}
}

type Operations []struct {
	Jitter               string
	Autonomous           int
	Group                string
	Chain                interface{}
	Use_learning_parsers bool
	Objective            interface{}
	Adversary            interface{}
	Auto_close           bool
	Visibility           int
	Name                 string `json:"name"`
	Id                   string `json:"id"`
	Obfuscator           string
	Host_group           []interface{}
	Planner              interface{}
	State                string
	Start                string
	Source               interface{}
}

// type Dummydata struct {
// 	UserId      string `json:"userId"`
// 	Status      string `json:"status"`
// 	ChallengeId string `json:"challengeId"`
// }

// var dummy []Dummydata = []Dummydata{{"user1", "good", "chall1"}, {"user1", "good", "chall1"}}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func connCheckAndSleep(conn *websocket.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	_, _, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		log.Println("read timeout")
		return false
	}
	time.Sleep(10 * time.Second)
	return true
}

func getReport(url string) OperationReport {
	b, _ := json.Marshal(struct {
		Enable_agent_output bool `json:"enable_agent_output"`
	}{Enable_agent_output: true})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	j := OperationReport{}
	json.NewDecoder(res.Body).Decode(&j)
	fmt.Println("Loop")
	return j
}

func handleWithId(conn *websocket.Conn, id string) {
	for {
		data := getReport(URL + "operations/" + id + "/report")
		if err := conn.WriteJSON(data); err != nil {
			log.Println(err)
			return
		}
		if !connCheckAndSleep(conn) {
			break
		}
	}
}

func getOperations() Operations {
	req, err := http.NewRequest("GET", "http://pdxf.tk:8888/api/v2/operations", nil)
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	j := Operations{}
	json.NewDecoder(res.Body).Decode(&j)
	fmt.Println("Loop alldashboard")
	return j
}

func handle(conn *websocket.Conn) {
	for {
		if err := conn.WriteJSON(struct {
			Data Operations `json:"data"`
		}{Data: getOperations()}); err != nil {
			log.Println(err)
			return
		}
		if !connCheckAndSleep(conn) {
			break
		}
	}
}

func SocketEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleError(err)
	defer conn.Close()
	fmt.Print("socket endpoint access\n")

	id := r.URL.Query().Get("id")

	if id == "" {
		handle(conn)
	} else {
		handleWithId(conn, id)
	}
}
