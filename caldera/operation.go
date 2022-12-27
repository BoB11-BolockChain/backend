package caldera

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/backend/utils"
)

// status -3 not executed

var OPERATION_API_URL = "http://www.pdxf.tk:8888/api/v2/operations"

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

func CreateOperation(name string, group string) string {
	createOperationBody := CreateOperationRequired{Name: name, Group: group, Auto_close: false}
	createOperationBody.Adversary.Adversary_iD = ""

	body, _ := json.Marshal(createOperationBody)
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
	return resOperation.Id
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

func AddPotentialLink(operationId string, linkRequired PotentialLinkBody) {
	url := OPERATION_API_URL + "/" + operationId + "/potential-links"

	body, _ := json.Marshal(linkRequired)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	utils.HandleError(err)
	req.Header.Add("KEY", "ADMIN123")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)

	if res.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(res.Body)
		utils.HandleError(err)
		log.Printf("potential link error %s\n", msg)
	}
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
