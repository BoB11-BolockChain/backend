package training

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/backend/database"
	_ "github.com/mattn/go-sqlite3"
)

type data struct {
	Scenario []Scenario `json:"scenario"`
}

type Scenario struct {
	Id			int			`json:"scene_id"`
	Title       string      `json:"scene_title"`
	Description string      `json:"scene_desc"`
	System      string      `json:"system"`
	Challenge   []Challenge `json:"challenge"`
}

type Challenge struct {
	Id          int    `json:"chall_id"`
	Title       string `json:"chall_title"`
	Description string `json:"chall_desc"`
	Score       string `json:"score"`
	Sequence    int    `json:"sequence"`
	Solved      string `json:"solved"`
}

type IdCheck struct {
	Id string
}

type FlagCheck struct {
	Id       string
	Chall_id int
	Flag     string
}

func Training(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	if r.Method == "POST" {
		var flag_check FlagCheck
		json.NewDecoder(r.Body).Decode(&flag_check)
		fmt.Println(flag_check.Chall_id)
		if flag_check.Chall_id == 0 {
			var data data

			scene_query := "SELECT id, title, description, system FROM scenario ORDER BY id LIMIT ?,1"

			row, err := database.DB().Query("SELECT COUNT(*) FROM scenario")
			if err != nil {
				panic(err)
			}
			defer row.Close()
			var scene_count int

			for row.Next() {
				row.Scan(&scene_count)
			}

			for i := 0; i < scene_count; i++ {
				var scene Scenario

				err := database.DB().QueryRow(scene_query, i).Scan(&scene.Id, &scene.Title, &scene.Description, &scene.System)
				if err != nil {
					panic(err)
				}
				data.Scenario = append(data.Scenario, scene)

				row, err := database.DB().Query("SELECT COUNT(*) FROM challenge WHERE scenario_id=?", scene.Id)
				if err != nil {
					panic(err)
				}
				defer row.Close()

				var chall_count int
				for row.Next() {
					row.Scan(&chall_count)
				}

				chall_query := "SELECT id, title, description, score, sequence FROM challenge WHERE scenario_id=? ORDER BY sequence LIMIT ?,1"

				chall_check_query := "SELECT id FROM challenge WHERE scenario_id=? AND sequence=?"

				chall_solved_check_query := "SELECT user_id FROM solved_challenge WHERE user_id=? AND solved_challenge_id=?"

				var chall_id int

				chall_solve_check := true

				for j := 0; j < chall_count; j++ {
					if !chall_solve_check {
						break
					}

					_ = database.DB().QueryRow(chall_check_query, scene.Id, j).Scan(&chall_id)
					// if chk_err != nil {
					// 	panic(chk_err)
					// }

					var dummy_id string
					solv_err := database.DB().QueryRow(chall_solved_check_query, flag_check.Id, chall_id).Scan(&dummy_id)
					if solv_err != nil && solv_err != sql.ErrNoRows {
						panic(solv_err)
					}

					var chall Challenge

					if dummy_id == "" {
						chall_solve_check = false
						chall.Solved = "False"
					} else {
						chall_solve_check = true
						chall.Solved = "True"
					}
					err := database.DB().QueryRow(chall_query, scene.Id, j).Scan(&chall.Id, &chall.Title, &chall.Description, &chall.Score, &chall.Sequence)
					if err != nil {
						panic(err)
					}
					data.Scenario[i].Challenge = append(data.Scenario[i].Challenge, chall)

					// if j != 0 {
					// 	chk_err := database.DB().QueryRow(chall_check_query, id, j-1).Scan(&chall_id)
					// 	if chk_err != nil {
					// 		panic(chk_err)
					// 	}

					// 	var dummy_id string
					// 	solv_err := database.DB().QueryRow(chall_solved_check_query, flag_check.Id, chall_id).Scan(&dummy_id)

					// 	if solv_err != nil && solv_err != sql.ErrNoRows {
					// 		panic(solv_err)
					// 	}

					// 	if dummy_id != "" {
					// 		var chall Challenge
					// 		err := database.DB().QueryRow(chall_query, id, j).Scan(&chall.Id, &chall.Title, &chall.Description, &chall.Score, &chall.Sequence)
					// 		if err != nil {
					// 			panic(err)
					// 		}
					// 		data.Scenario[i].Challenge = append(data.Scenario[i].Challenge, chall)
					// 	} else {
					// 		continue
					// 	}
					// } else {
					// 	var chall Challenge
					// 	err := database.DB().QueryRow(chall_query, id, j).Scan(&chall.Id, &chall.Title, &chall.Description, &chall.Score, &chall.Sequence)
					// 	if err != nil {
					// 		panic(err)
					// 	}
					// 	data.Scenario[i].Challenge = append(data.Scenario[i].Challenge, chall)
					// }
				}

			}
			json.NewEncoder(w).Encode(data)
			fmt.Println(data)
		} else {
			var flag string
			fmt.Println(flag_check)
			flag_query := "SELECT flag FROM challenge WHERE id=? AND flag=?"
			err := database.DB().QueryRow(flag_query, flag_check.Chall_id, flag_check.Flag).Scan(&flag)
			fmt.Println(flag)
			var solve_check string
			if err != nil {
				solve_check = "False"
			} else {
				var dummy_id string
				fmt.Println(flag_check.Id)
				fmt.Println(flag_check.Chall_id)
				flag_select_query := "SELECT user_id FROM solved_challenge WHERE user_id=? AND solved_challenge_id=?"
				err := database.DB().QueryRow(flag_select_query, flag_check.Id, flag_check.Chall_id).Scan(&dummy_id)
				fmt.Println("flag 값?")
				fmt.Println(dummy_id)
				if err != nil && err != sql.ErrNoRows {
					panic(err)
				}

				if dummy_id == "" {
					now := time.Now()
					now_parsing := now.Format("2006-01-02 15:04:05")

					insert, _ := database.DB().Prepare("INSERT INTO solved_challenge (user_id, solved_challenge_id, solved_time) VALUES (?, ?, ?)")
					_, err := insert.Exec(flag_check.Id, flag_check.Chall_id, now_parsing)
					if err != nil {
						panic(err)
					}
					solve_check = "Solve a Challenge"
				} else {
					solve_check = "Aleady Solved"
				}
			}
			fmt.Println(solve_check)
			check := struct {
				Chall_id string `json:"chall_id"`
			}{solve_check}
			json.NewEncoder(w).Encode(check)
		}
	}
}
