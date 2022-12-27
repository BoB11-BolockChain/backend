package training

import (
	"encoding/json"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

type TrainingHash struct {
	Id          int             `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	System      string          `json:"system"`
	Score       int             `json:"score"`
	Vm_name     string          `json:"vm_name"`
	Vm_id       string          `json:"vm_id"`
	Vm_pw       string          `json:"vm_pw"`
	Visible     bool            `json:"visible"`
	Challenges  []ChallengeHash `json:"challenges"`
}
type ChallengeHash struct {
	Id          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Score       int          `json:"score"`
	Flag        string       `json:"flag"`
	Sequence    int          `json:"sequence"`
	Tactics     []TacticHash `json:"tactics"`
	Hash        string       `json:"hash"`
}
type TacticHash struct {
	Id       int           `json:"id"`
	Title    string        `json:"title"`
	Sequence int           `json:"sequence"`
	Delay    int           `json:"delay"`
	Payloads []PayloadHash `json:"payloads"`
	Hash     string        `json:"hash"`
}
type PayloadHash struct {
	Id       int    `json:"id"`
	Payload  string `json:"payload"`
	Sequence int    `json:"sequence"`
	Hash     string `json:"hash"`
}

func GetTraining(w http.ResponseWriter, r *http.Request) {
	trainingId := r.URL.Query().Get("trainingId")
	db := database.DB()

	t := TrainingHash{}
	err := db.QueryRow("select id,title,description,system,score,vm_name,vm_id,vm_pw,visible from scenario where id=?", trainingId).Scan(&t.Id, &t.Title, &t.Description, &t.System, &t.Score, &t.Vm_name, &t.Vm_id, &t.Vm_pw, &t.Visible)
	utils.HandleError(err)

	chRows, err := db.Query("select id,title,description,score,flag,sequence from challenge where scenario_id=?", trainingId)
	utils.HandleError(err)
	defer chRows.Close()

	for chRows.Next() {
		c := ChallengeHash{}
		err = chRows.Scan(&c.Id, &c.Title, &c.Description, &c.Score, &c.Flag, &c.Sequence)
		utils.HandleError(err)
		c.Hash = utils.Hash(c)

		rows, err := db.Query("select id,title,sequence,delay from tactic where challenge_id=?", c.Id)
		utils.HandleError(err)
		defer rows.Close()

		for rows.Next() {
			ta := TacticHash{}
			err = rows.Scan(&ta.Id, &ta.Title, &ta.Sequence, &ta.Delay)
			utils.HandleError(err)
			ta.Hash = utils.Hash(ta)

			payloads, err := db.Query("select id,payload,sequence from payload where tactic_id=?", ta.Id)
			utils.HandleError(err)
			defer payloads.Close()

			for payloads.Next() {
				p := PayloadHash{}
				err = payloads.Scan(&p.Id, &p.Payload, &p.Sequence)
				utils.HandleError(err)
				p.Hash = utils.Hash(p)

				ta.Payloads = append(ta.Payloads, p)
			}
			c.Tactics = append(c.Tactics, ta)
		}
		t.Challenges = append(t.Challenges, c)
	}

	json.NewEncoder(w).Encode(t)
}

func GetAllTrainings(w http.ResponseWriter, r *http.Request) {
	db := database.DB()
	rows, err := db.Query("select id,title,description,system,score,visible from scenario")
	utils.HandleError(err)

	trainings := make([]TrainingHash, 0)
	for rows.Next() {
		s := TrainingHash{}
		rows.Scan(&s.Id, &s.Title, &s.Description, &s.System, &s.Score, &s.Visible)
		trainings = append(trainings, s)
	}

	json.NewEncoder(w).Encode(trainings)
}
