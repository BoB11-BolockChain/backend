package dashboard

import (
	"encoding/json"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

type SolvedScenario struct {
	Id         int               `json:"id"`
	Title      string            `json:"title"`
	Challenges []SolvedChallenge `json:"challenges"`
}

type SolvedChallenge struct {
	Id     int    `json:"id"`
	Solved bool   `json:"solved"`
	Score  int    `json:"score"`
	Title  string `json:"title"`
}

type DashboardUser struct {
	Id             string `json:"id"`
	LastSolution   string `json:"lastSolution"`
	LastConnection string `json:"lastConnection"`
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	db := database.DB()
	rows, err := db.Query("select id from user")
	utils.HandleError(err)

	users := make([]DashboardUser, 0)

	for rows.Next() {
		user := DashboardUser{}
		rows.Scan(&user.Id)

		r := db.QueryRow("select c.title from solved_challenge s inner join challenge c on s.solved_challenge_id=c.id where s.user_id=? order by s.solved_time desc;", user.Id)
		r.Scan(&user.LastSolution)
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

// get scenarios and its challenges with solved or not
func DashboardByUser(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	db := database.DB()

	rows, err := db.Query("select s.solved_challenge_id from solved_challenge s inner join user u on s.user_id=u.id where u.id=?", userId)
	utils.HandleError(err)

	solvedIds := make(map[int]int, 0)
	for rows.Next() {
		var scid int
		err = rows.Scan(&scid)
		utils.HandleError(err)
		solvedIds[scid] = scid
	}

	rows, err = db.Query("select id, title from scenario")
	utils.HandleError(err)

	solvedScenarios := make([]SolvedScenario, 0)

	for rows.Next() {
		scenario := SolvedScenario{}
		err = rows.Scan(&scenario.Id, &scenario.Title)
		utils.HandleError(err)

		cRows, err := db.Query("select id,title,score from challenge where scenario_id=?", scenario.Id)
		utils.HandleError(err)

		for cRows.Next() {
			challenge := SolvedChallenge{}
			err = cRows.Scan(&challenge.Id, &challenge.Title, &challenge.Score)
			utils.HandleError(err)

			_, ok := solvedIds[challenge.Id]
			if ok {
				challenge.Solved = true
			}
			scenario.Challenges = append(scenario.Challenges, challenge)
		}
		// incident response result here
		// table with operationid required
		solvedScenarios = append(solvedScenarios, scenario)
	}
	json.NewEncoder(w).Encode(solvedScenarios)
}
