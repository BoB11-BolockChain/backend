package training

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	_ "github.com/mattn/go-sqlite3"
)

type ChallCheck struct {
	Id             string
	Chall_id       int
	Chall_sequence int
}

type ChallState struct {
	Solvestate string `json:"solve_state"` // 풀었는지 안풀었는지
}

func ChallengeCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	if r.Method == "POST" {
		var ch_ck ChallCheck
		var ch_st ChallState
		json.NewDecoder(r.Body).Decode(&ch_ck)
		fmt.Println(ch_ck)
		var ch_id int

		chall_solve_query := "SELECT solved_challenge_id FROM solved_challenge WHERE user_id=? AND solved_challenge_id=?"
		err := database.DB().QueryRow(chall_solve_query, ch_ck.Id, ch_ck.Chall_id).Scan(&ch_id)
		if err != nil {
			ch_st.Solvestate = "False"
		} else {
			ch_st.Solvestate = "True"
		}

		json.NewEncoder(w).Encode(ch_st)
	}
}
