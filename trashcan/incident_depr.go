package trashcan

// package caldera

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/backend/database"
// 	"github.com/backend/utils"
// 	"github.com/gorilla/websocket"
// )

// // caldera data. should be replaced with challenge db
// const URL = "http://pdxf.tk:8888/api/v2/"

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin:     func(r *http.Request) bool { return true },
// }

// func connCheckAndSleep(conn *websocket.Conn) bool {
// 	conn.SetReadDeadline(time.Now().Add(15 * time.Second))
// 	_, _, err := conn.ReadMessage()
// 	if err != nil {
// 		conn.Close()
// 		log.Println("read timeout")
// 		return false
// 	}
// 	time.Sleep(10 * time.Second)
// 	return true
// }

// func getReport(url string) OperationReport {
// 	b, _ := json.Marshal(struct {
// 		Enable_agent_output bool `json:"enable_agent_output"`
// 	}{Enable_agent_output: true})
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
// 	utils.HandleError(err)

// 	req.Header.Add("KEY", "ADMIN123")

// 	client := &http.Client{}
// 	res, err := client.Do(req)
// 	utils.HandleError(err)
// 	defer res.Body.Close()

// 	j := OperationReport{}
// 	json.NewDecoder(res.Body).Decode(&j)
// 	fmt.Println("Loop byid")
// 	return j
// }

// func handleWithId(conn *websocket.Conn, id string) {
// 	for {
// 		data := getReport(URL + "operations/" + id + "/report")
// 		if err := conn.WriteJSON(data); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		if !connCheckAndSleep(conn) {
// 			break
// 		}
// 	}
// }

// func getOperations() []Operation {
// 	req, err := http.NewRequest("GET", "http://pdxf.tk:8888/api/v2/operations", nil)
// 	utils.HandleError(err)

// 	req.Header.Add("KEY", "ADMIN123")

// 	client := &http.Client{}
// 	res, err := client.Do(req)
// 	utils.HandleError(err)
// 	defer res.Body.Close()

// 	j := make([]Operation, 0)
// 	json.NewDecoder(res.Body).Decode(&j)
// 	fmt.Println("Loop alldashboard")
// 	return j
// }

// func handle(conn *websocket.Conn) {
// 	for {
// 		if err := conn.WriteJSON(struct {
// 			Data []Operation `json:"data"`
// 		}{Data: getOperations()}); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		if !connCheckAndSleep(conn) {
// 			break
// 		}
// 	}
// }

// func SocketEndpoint(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	utils.HandleError(err)
// 	defer conn.Close()
// 	fmt.Print("socket endpoint access\n")

// 	id := r.URL.Query().Get("id")

// 	if id == "" {
// 		handle(conn)
// 	} else {
// 		handleWithId(conn, id)
// 	}
// }

// func getOperation(operationId string) Operation {
// 	url := OPERATION_API_URL + "/" + operationId
// 	req, err := http.NewRequest("GET", url, nil)
// 	utils.HandleError(err)

// 	req.Header.Add("KEY", "ADMIN123")

// 	client := &http.Client{}
// 	res, err := client.Do(req)
// 	utils.HandleError(err)
// 	defer res.Body.Close()

// 	result := Operation{}
// 	json.NewDecoder(res.Body).Decode(&result)
// 	return result
// }

// func handleOperations(conn *websocket.Conn, userId string, scenarioId int) {
// 	db := database.DB()

// 	// assume incident response is one
// 	row := db.QueryRow("select operation_id from solved_scenario where user_id=? and solved_scenario_id=?", userId, scenarioId)
// 	var operationId string
// 	row.Scan(&operationId)

// 	for connCheckAndSleep(conn) {
// 		operation := getOperation(operationId)
// 		conn.WriteJSON(operation)
// 		if operation.State == "finished" {
// 			conn.Close()
// 			break
// 		}
// 	}
// }
