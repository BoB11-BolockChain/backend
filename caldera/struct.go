package caldera

import "time"

// OperationReport.Steps {[paw]:[]Steps, [paw]:[]Steps}
type OperationReport struct {
	Name       string        `json:"name"`
	Start      string        `json:"start"`
	Host_group []interface{} `json:"host_group"`
	Steps      interface{}   `json:"steps"`
	Finish     bool          `json:"finish"`
	Planner    string
	Adversary  interface{}
	Jitter     string
	Objectives interface{}
	Facts      []interface{}
}

type Step struct {
	LinkID           string    `json:"link_id"`
	AbilityID        string    `json:"ability_id"`
	Command          string    `json:"command"`
	PlaintextCommand string    `json:"plaintext_command"`
	Delegated        time.Time `json:"delegated"`
	Run              time.Time `json:"run"`
	Status           int       `json:"status"`
	Platform         string    `json:"platform"`
	Executor         string    `json:"executor"`
	Pid              int       `json:"pid"`
	Description      string    `json:"description"`
	Name             string    `json:"name"`
	Attack           struct {
		Tactic        string `json:"tactic"`
		TechniqueName string `json:"technique_name"`
		TechniqueID   string `json:"technique_id"`
	} `json:"attack"`
	Output string `json:"output"`
}
