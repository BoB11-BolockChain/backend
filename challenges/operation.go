package challenges

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/utils"
)

// v4.0.0
type OperationReqBody struct {
	Jitter               string
	Autonomous           int
	Group                string
	Use_learning_parsers bool
	Objective            interface{}
	Adversary            interface{}
	Auto_close           bool
	Visibility           int
	Name                 string
	Id                   string
	Obfuscator           string
	Planner              interface{}
	State                string
	Source               interface{}
}

type RequiredFields struct {
	Name      string `json:"name"`
	Adversary struct {
		AdversaryID string `json:"adversary_id"`
	} `json:"adversary"`
	Planner struct {
		ID string `json:"id"`
	} `json:"planner"`
	Source struct {
		ID string `json:"id"`
	} `json:"source"`
	Group string `json:"group"`
	State string `json:"state"`
}

var operationURL = "http://www.pdxf.tk:8888/api/v2/operations"

func CreateOperation(w http.ResponseWriter, r *http.Request) {
	required := RequiredFields{}
	json.NewDecoder(r.Body).Decode(&required)
	body, _ := json.Marshal(required)
	req, err := http.NewRequest("POST", operationURL, bytes.NewBuffer(body))
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	resOperation := struct {
		Id string `json:"id"`
	}{}
	json.NewDecoder(res.Body).Decode(&resOperation)

	fmt.Print("operation created ", resOperation.Id)

	json.NewEncoder(w).Encode(resOperation)
}
