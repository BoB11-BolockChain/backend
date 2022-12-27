package training

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

func CreateTraining(w http.ResponseWriter, r *http.Request) {
	create := TrainingHash{}
	err := json.NewDecoder(r.Body).Decode(&create)
	utils.HandleError(err)

	db := database.DB()

	var visible int
	if create.Visible {
		visible = 1
	}
	res, err := db.Exec("insert into scenario(title,description,system,score,vm_name,vm_id,vm_pw,visible) values(?,?,?,?,?,?,?,?)", create.Title, create.Description, create.System, create.Score, create.Vm_name, create.Vm_id, create.Vm_pw, visible)
	utils.HandleError(err)

	scenarioId, _ := res.LastInsertId()

	stmtc, err := db.Prepare("insert into challenge(title,description,score,flag,sequence,scenario_id) values(?,?,?,?,?,?)")
	utils.HandleError(err)
	defer stmtc.Close()

	stmtt, err := db.Prepare("insert into tactic(title,sequence,challenge_id,delay) values(?,?,?,?)")
	utils.HandleError(err)
	defer stmtt.Close()

	stmtp, err := db.Prepare("insert into payload(payload,sequence,tactic_id) values(?,?,?)")
	utils.HandleError(err)
	defer stmtp.Close()

	for i, v := range create.Challenges {
		cres, err := stmtc.Exec(v.Title, v.Description, v.Score, v.Flag, i+1, scenarioId)
		utils.HandleError(err)

		challengeId, _ := cres.LastInsertId()

		for j, t := range v.Tactics {
			tres, err := stmtt.Exec(t.Title, j+1, challengeId, t.Delay)
			utils.HandleError(err)

			tacticId, _ := tres.LastInsertId()

			for k, p := range t.Payloads {
				_, err = stmtp.Exec(p.Payload, k+1, tacticId)
				utils.HandleError(err)
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func EditTraining(w http.ResponseWriter, r *http.Request) {
	body := TrainingHash{}
	err := json.NewDecoder(r.Body).Decode(&body)
	utils.HandleError(err)

	db := database.DB()
	db.Exec("PRAGMA foreign_keys = ON")

	tx, err := db.Begin()
	utils.HandleError(err)
	defer tx.Rollback()

	tx.Exec("PRAGMA foreign_keys = ON")

	visible := 0
	if body.Visible {
		visible = 1
	}
	if body.Id == 0 {
		_, nerr := tx.Exec("insert into scenario(title,description,system,score,vm_name,vm_id,vm_pw,visible) values(?,?,?,?,?,?,?,?)", body.Title, body.Description, body.System, body.Score, body.Vm_name, body.Vm_id, body.Vm_pw, visible)
		utils.HandleError(nerr)
	} else {
		_, nerr := tx.Exec("update scenario set title=?, description=?, system=?, score=?, vm_name=?,vm_id=?,vm_pw=?,visible=? where id=?", body.Title, body.Description, body.System, body.Score, body.Vm_name, body.Vm_id, body.Vm_pw, visible, body.Id)
		utils.HandleError(nerr)
	}

	_, err = tx.Exec("delete from challenge where scenario_id=?", body.Id)
	utils.HandleError(err)

	for i, v := range body.Challenges {
		res, err := tx.Exec("insert into challenge(title,description,score,flag,sequence,scenario_id) values(?,?,?,?,?,?)", v.Title, v.Description, v.Score, v.Flag, i+1, body.Id)
		utils.HandleError(err)

		challCreatedId, err := res.LastInsertId()
		utils.HandleError(err)

		for j, t := range v.Tactics {
			tres, err := tx.Exec("insert into tactic(title,sequence,challenge_id,delay) values(?,?,?,?)", t.Title, j+1, challCreatedId, t.Delay)
			utils.HandleError(err)

			tacticCreatedId, err := tres.LastInsertId()
			utils.HandleError(err)

			for k, p := range t.Payloads {
				_, err := tx.Exec("insert into payload(payload,sequence,tactic_id) values(?,?,?)", p.Payload, k+1, tacticCreatedId)
				utils.HandleError(err)
			}
		}
	}

	err = tx.Commit()
	utils.HandleError(err)
	w.WriteHeader(http.StatusAccepted)
}

func DeleteTraining(w http.ResponseWriter, r *http.Request) {
	trainingId := r.URL.Query().Get("trainingId")
	db := database.DB()

	db.Exec("PRAGMA foreign_keys = ON")

	res, err := db.Exec("delete from scenario where id=?", trainingId)
	utils.HandleError(err)

	fmt.Print(res.RowsAffected())
	w.WriteHeader(http.StatusAccepted)
}
