package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
	_ "github.com/mattn/go-sqlite3"
)

// var ability []ability

// branch := struct {
// 	Branch []ability `json:"branch"`
// }{ability}

// data := struct {
// 	Data []branch `json:"data"`
// }{branch}

// type Ability struct {
// 	Data []struct {
// 		Branch      string `json:"branch"`
// 		Seq         string `json:"seq"`
// 		Payload     string
// 		AbilityName string
// 	} `json:"data"`
// }

// type ChallengeNum struct {
// 	Num string
// }

// type SendAbility struct {
// 	Payload     string
// 	AbilityName string
// }

// type SendBranch struct {
// 	Data []struct {
// 		Payload     string
// 		AbilityName string
// 	}
// }

// func LoadBasic(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	// apt := "1"
// 	var data SendAbility
// 	var datas []SendAbility
// 	query := "select payload, abilityname from basic where name=1 ORDER BY branch, seq"
// 	print(query)
// 	rows, err := database.DB().Query(query)
// 	print(rows)

// 	fmt.Println(err)
// 	defer rows.Close()

// 	for rows.Next() {
// 		rows.Scan(&data.Payload, &data.AbilityName)
// 		fmt.Println(data)
// 		datas = append(datas, data)
// 		fmt.Println(datas)
// 	}
// }

type data struct {
	Data []struct {
		Branch []struct {
			Seq         int    `json:"Sequence"`
			Payload     string `json:"Payload"`
			Abilityname string `json:"AbilityName"`
		} `json:"Branch"`
		//Delete string `json:Delete`
	} `json:"data"`
	Info struct {
		Title       string `json:"Title"`
		Description string `json:"Description"`
		Score       int    `json:"Score"`
		Os          string `json:"OS"`
	} `json:"info"`
}

func InsertData2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	if r.Method == "POST" {
		var ch_dat data

		json.NewDecoder(r.Body).Decode(&ch_dat)
		fmt.Println(ch_dat)
		fmt.Println(ch_dat.Data)
		fmt.Println(ch_dat.Data[0].Branch[0].Payload)
		// abil_insert, _ := database.DB().Prepare("INSERT INTO ability (num, abil_name, payload, branch_num, seq) VALUES ()")
		// branch_insert, _ := database.DB().Prepare("INSERT INTO branch (num, chall_num, seq, ability_cnt) VALUES ()")
		chall_insert, _ := database.DB().Prepare("INSERT INTO challenge (num, title, desc, os, score) VALUES (NULL, ?, ?, ?, ?)")
		_, err := chall_insert.Exec(ch_dat.Info.Title, ch_dat.Info.Description, ch_dat.Info.Os, ch_dat.Info.Score)
		if err != nil {
			utils.HandleError(err)
		}
		var chall_num int
		call_num_query := "SELECT num FROM challenge WHERE title=? ORDER BY rowid DESC LIMIT 1"
		row := database.DB().QueryRow(call_num_query, ch_dat.Info.Title)
		row.Scan(&chall_num)

		for i, branch := range ch_dat.Data {
			branch_insert, _ := database.DB().Prepare("INSERT INTO branch (num, chall_num, seq) VALUES (NULL, ?, ?)")
			_, err := branch_insert.Exec(chall_num, i+1)
			if err != nil {
				utils.HandleError(err)
			}

			var branch_num int
			branch_num_query := "SELECT num FROM branch WHERE chall_num=? AND seq=?"
			row := database.DB().QueryRow(branch_num_query, chall_num, i+1)
			row.Scan(&branch_num)

			for j, ability := range branch.Branch {
				abil_insert, _ := database.DB().Prepare("INSERT INTO ability (num, abil_name, payload, branch_num, seq) VALUES (NULL, ?, ?, ?, ?)")
				_, err := abil_insert.Exec(ability.Abilityname, ability.Payload, branch_num, j)
				if err != nil {
					utils.HandleError(err)
				}
			}
		}
	}
}

// func PrintData(w http.ResponseWriter, r *http.Request) {
// 	// w.Header().Set("Access-Control-Allow-Origin", "*")
// 	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 	r.ParseForm()
// 	var ch_num ChallengeNum
// 	var data SendAbility
// 	json.NewDecoder(r.Body).Decode(&ch_num)
// 	fmt.Println(ch_num)

// 	// get branch_num (1, 2, 3)
// 	query := fmt.Sprintf("select payload, abilityname from branch  where ch_num=%s ORDER BY branch_num, seq", ch_num)
// 	print(query)
// 	rows, err := database.DB().Query(query)
// 	fmt.Println(err)
// 	defer rows.Close()

// 	for rows.Next() {
// 		rows.Scan(&data.Payload, &data.AbilityName)
// 		fmt.Println(data)
// 	}

// 	// query := fmt.Sprintf("select payload, abilityname from ability where ='%s'", session.SessionId)

// }
