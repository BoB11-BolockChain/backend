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

// Required: "name", "adversary.adversary_id", "planner.planner_id", and "source.id"
type RequiredFields struct {
	Name         string `json:"name"`
	Adversary_id string `json:"adversary_id"`
	Planner_id   string `json:"planner_id"`
	Source_id    string `json:"source_id"`
	Group        string `json:"group"`
}

var operationURL = "http://pdxf.tk:8888/api/v2/operations"

func CreateOperation(w http.ResponseWriter, r *http.Request) {
	required := RequiredFields{}
	json.NewDecoder(r.Body).Decode(&required)
	fmt.Print(required)

	body, _ := json.Marshal(required)
	req, err := http.NewRequest("POST", operationURL, bytes.NewBuffer(body))
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")

	client := &http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	resOperation := struct {
		Id string `json:"id"`
	}{}
	json.NewDecoder(res.Body).Decode(&resOperation)

	fmt.Print(resOperation)

	json.NewEncoder(w).Encode(resOperation)
}
