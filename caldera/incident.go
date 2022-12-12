package caldera

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/backend/database"
	"github.com/backend/utils"
	"github.com/gorilla/websocket"
)

type Operation struct {
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func SocketEndpoint(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	utils.HandleError(err)
	defer conn.Close()

	userId := r.URL.Query().Get("userId")
	scenarioId := r.URL.Query().Get("scenarioId")
	db := database.DB()
	row := db.QueryRow("select operation_id from solved_scenario where user_id=? and solved_scenario_id=?", userId, scenarioId)

	var operationId string
	err = row.Scan(&operationId)
	utils.HandleError(err)

	for {
		data := GetOperation(operationId)

		if err := conn.WriteJSON(data); err != nil {
			log.Println(err)
			break
		}

		if data.State == "finished" {
			log.Print("operation finished.")
			break
		}

		if !connAliveCheck(conn) {
			log.Print("connection closed : read timeout")
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func connAliveCheck(conn *websocket.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
	_, _, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		log.Println("read timeout")
		return false
	}
	return true
}

func GetOperation(operationId string) Operation {
	url := OPERATION_API_URL + "/" + operationId
	req, err := http.NewRequest("GET", url, nil)
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	result := Operation{}
	json.NewDecoder(res.Body).Decode(&result)
	return result
}
