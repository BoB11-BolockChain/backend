package challenges

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/backend/utils"
)

// v4.0.0
type OperationReqBody struct {
	jitter               string
	autonomous           int
	group                string
	use_learning_parsers bool
	objective            interface{}
	adversary            interface{}
	auto_close           bool
	visibility           int
	name                 string
	id                   string
	obfuscator           string
	planner              interface{}
	state                string
	source               interface{}
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

	json.NewEncoder(w).Encode(resOperation)
}
