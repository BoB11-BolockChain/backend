package training

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

type CreateTrainingDTO struct {
	Id         int            `json:"id"`
	Scenario   ScenarioDTO    `json:"scenario"`
	Challenges []ChallengeDTO `json:"challenges"`
}
type ScenarioDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	System      string `json:"system"`
	Score       int    `json:"score"`
}
type ChallengeDTO struct {
	Id          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Score       int         `json:"score"`
	Flag        string      `json:"flag"`
	Sequence    int         `json:"sequence"`
	Tactics     []TacticDTO `json:"tactics"`
}
type TacticDTO struct {
	Id       int          `json:"id"`
	Title    string       `json:"title"`
	Sequence int          `json:"sequence"`
	Payloads []PayloadDTO `json:"payloads"`
}
type PayloadDTO struct {
	Id       int    `json:"id"`
	Payload  string `json:"payload"`
	Sequence int    `json:"sequence"`
}

func CreateTraining(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("r.Body: %v\n", r.Body)
	create_dto := CreateTrainingDTO{}
	err := json.NewDecoder(r.Body).Decode(&create_dto)
	utils.HandleError(err)

	db := database.DB()

	args := utils.GetStructValues(create_dto.Scenario)
	res, err := db.Exec("insert into scenario(title,description,system,score) values(?,?,?,?)", args...)
	utils.HandleError(err)

	scenarioId, _ := res.LastInsertId()

	stmtc, err := db.Prepare("insert into challenge(title,description,score,flag,sequence,scenario_id) values(?,?,?,?,?,?)")
	utils.HandleError(err)
	defer stmtc.Close()

	stmtt, err := db.Prepare("insert into tactic(title, sequence, challenge_id) values(?,?,?)")
	utils.HandleError(err)
	defer stmtt.Close()

	stmtp, err := db.Prepare("insert into payload(payload,sequence,tactic_id) values(?,?,?)")
	utils.HandleError(err)
	defer stmtp.Close()

	for i, v := range create_dto.Challenges {
		cres, err := stmtc.Exec(v.Title, v.Description, v.Score, v.Flag, i+1, scenarioId)
		utils.HandleError(err)

		challengeId, _ := cres.LastInsertId()

		for j, t := range v.Tactics {
			tres, err := stmtt.Exec(t.Title, j+1, challengeId)
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

	tx, err := database.DB().Begin()
	utils.HandleError(err)
	defer tx.Rollback()

	tx.Exec("PRAGMA foreign_keys = ON")

	if body.Id == 0 {
		_, nerr := tx.Exec("insert into scenario(title,description,system,score) values(?,?,?,?)", body.Title, body.Description, body.System, body.Score)
		utils.HandleError(nerr)
	} else {
		_, nerr := tx.Exec("update scenario set title=?, description=?, system=?, score=? where id=?", body.Title, body.Description, body.System, body.Score, body.Id)
		utils.HandleError(nerr)
	}

	//remove chall from db if not exists
	cMap := make(map[int]ChallengeHash)
	for _, ch := range body.Challenges {
		cMap[ch.Id] = ch
	}
	rows, err := tx.Query("select id from challenge where scenario_id=?", body.Id)
	utils.HandleError(err)
	for rows.Next() {
		var challId int
		rows.Scan(&challId)
		_, a := cMap[challId]
		if !a {
			_, err = tx.Exec("delete from challenge where id=?", challId)
			utils.HandleError(err)
		}
	}

	for i, ch := range body.Challenges {
		var challId int64
		// ch.id can be 0(nil) and something exists

		var dummy int
		err = tx.QueryRow("select id from challenge where id=?", ch.Id).Scan(&dummy)
		if err == sql.ErrNoRows {
			res, err := tx.Exec("insert into challenge(title,description,score,flag,sequence,scenario_id) values(?,?,?,?,?,?)", ch.Title, ch.Description, ch.Score, ch.Flag, i+1, body.Id)
			utils.HandleError(err)
			challId, err = res.LastInsertId()
			utils.HandleError(err)
		} else {
			utils.HandleError(err)
			challId = int64(ch.Id)
			_, err = tx.Exec("update challenge set title=?,description=?,score=?,flag=?,sequence=? where id=?", ch.Title, ch.Description, ch.Score, ch.Flag, i+1, challId)
			utils.HandleError(err)
			_, err = tx.Exec("delete from tactic where challenge_id=?", challId)
			utils.HandleError(err)
		}

		for i2, th := range ch.Tactics {
			var tacticId int64
			res, err := tx.Exec("insert into tactic(title, sequence, challenge_id) values(?,?,?)", th.Title, i2+1, challId)
			utils.HandleError(err)
			tacticId, err = res.LastInsertId()
			utils.HandleError(err)

			for i3, ph := range th.Payloads {
				_, err = tx.Exec("insert into payload(payload,sequence,tactic_id) values(?,?,?)", ph.Payload, i3+1, tacticId)
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
