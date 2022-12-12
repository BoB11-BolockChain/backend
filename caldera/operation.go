package caldera

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

// status -3 not executed

var OPERATION_API_URL = "http://www.pdxf.tk:8888/api/v2/operations"

type CreateOperationBody struct {
	IsIR       bool
	UserId     string
	ScenarioId string
	required   CreateOperationRequired
}

type CreateOperationRequired struct {
	Name      string `json:"name"`
	Adversary struct {
		Adversary_iD string `json:"adversary_id"`
	} `json:"adversary"`
	Planner struct {
		ID string `json:"id"`
	} `json:"planner"`
	Group      string `json:"group"`
	Auto_close bool   `json:"auto_close"`
}

func CreateOperation(w http.ResponseWriter, r *http.Request) {
	createOperationBody := CreateOperationBody{}
	json.NewDecoder(r.Body).Decode(&createOperationBody)

	body, _ := json.Marshal(createOperationBody.required)
	req, err := http.NewRequest("POST", OPERATION_API_URL, bytes.NewBuffer(body))
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	resOperation := struct {
		Id string `json:"id"`
	}{}
	json.NewDecoder(res.Body).Decode(&resOperation)

	if createOperationBody.IsIR {
		db := database.DB()
		_, err := db.Exec("insert into solved_scenario(user_id,scenario_id,operation_id) values(?,?,?)", createOperationBody.UserId, createOperationBody.ScenarioId, resOperation.Id)
		utils.HandleError(err)
	}

	json.NewEncoder(w).Encode(resOperation)
}

type PotentialLinkBody struct {
	Paw     string `json:"paw"`
	Ability struct {
		Name string `json:"name"`
	} `json:"ability"`
	Executor struct {
		Platform string `json:"platform"`
		Name     string `json:"name"`
		Command  string `json:"command"`
	} `json:"executor"`
}

func AddPotentialLink(w http.ResponseWriter, r *http.Request) {
	operationId := r.URL.Query().Get("operationId")
	potentialLinkBody := PotentialLinkBody{}
	err := json.NewDecoder(r.Body).Decode(&potentialLinkBody)
	utils.HandleError(err)

	url := OPERATION_API_URL + "/" + operationId + "/potential-links"

	body, _ := json.Marshal(potentialLinkBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	utils.HandleError(err)
	req.Header.Add("KEY", "ADMIN123")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)

	if res.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func TerminateOperation(operationId string) bool {
	state := struct {
		State string `json:"state"`
	}{State: "cleanup"}
	body, err := json.Marshal(state)
	utils.HandleError(err)

	url := OPERATION_API_URL + "/" + operationId
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	utils.HandleError(err)

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)

	json.NewDecoder(res.Body).Decode(&state)
	return state.State == "cleanup"
}

func StartScenario(scenarioId int, userId string, operationId string) {
	db := database.DB()

	c_ids, err := db.Query("select id from challenge where scenario_id=?", scenarioId)
	utils.HandleError(err)

	for c_ids.Next() {
		var challengeId int
		err = c_ids.Scan(&challengeId)
		utils.HandleError(err)

		// just execute payloads simultaneously
		payloads, err := db.Query("select p.payload from tactic t inner join payload p on t.id=p.tactic_id where t.challenge_id=? order by t.sequence,p.sequence", challengeId)
		utils.HandleError(err)
		for payloads.Next() {
			var payload string
			err = payloads.Scan(&payload)
			utils.HandleError(err)

			potentialLinkBody := PotentialLinkBody{}
			potentialLinkBody.Paw = userId
			potentialLinkBody.Ability.Name = ""
			potentialLinkBody.Executor.Name = "cmd"
			potentialLinkBody.Executor.Platform = "windows"
			potentialLinkBody.Executor.Command = payload

			url := OPERATION_API_URL + "/" + operationId + "/potential-links"

			body, _ := json.Marshal(potentialLinkBody)
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
			utils.HandleError(err)
			req.Header.Add("KEY", "ADMIN123")
			req.Header.Add("Content-Type", "application/json; charset=utf-8")

			client := &http.Client{}
			_, err = client.Do(req)
			utils.HandleError(err)
		}
	}
}
