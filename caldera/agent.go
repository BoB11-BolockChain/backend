package caldera

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/backend/utils"
	"golang.org/x/exp/slices"
)

func IsAgentAlive(paw string, executorName string) bool {
	req, err := http.NewRequest("GET", "http://pdxf.tk:8888/api/v2/agents/"+paw, nil)
	utils.HandleError(err)
	req.Header.Add("KEY", "ADMIN123")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)

	if res.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(res.Body)
		utils.HandleError(err)
		log.Printf("Caldera API Error : GET agents : %s\n", msg)
		return false
	}

	resBody := struct {
		Executors []string `json:"executors"`
		Last_seen string   `json:"last_seen"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	utils.HandleError(err)

	gmtTime := time.Now().In(time.FixedZone("GMT+0", 0))
	lastSeenTime, _ := time.Parse(time.RFC3339, resBody.Last_seen)
	if lastSeenTime.Add(time.Minute).Before(gmtTime) {
		log.Printf("AgentAliveCheck : Agent %s is dead.\n", paw)
		return false
	}

	if !slices.Contains(resBody.Executors, executorName) {
		log.Printf("AgentAliveCheck : Executor %s not found in %s\n", executorName, paw)
		return false
	}
	return true
}
