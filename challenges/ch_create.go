package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
	_ "github.com/mattn/go-sqlite3"
)

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

type Send struct {
	Payload     string `json:"Payload"`
	Abilityname string `json:"AbilityName"`
}

type Sendbranch struct {
	BranchName []Send
}

func InsertData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	if r.Method == "POST" {
		var ch_dat data

		json.NewDecoder(r.Body).Decode(&ch_dat)
		fmt.Println(ch_dat)
		fmt.Println(ch_dat.Data)
		fmt.Println(ch_dat.Data[0].Branch[0].Payload)
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

func LoadBasic(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// apt := "1"

	var b_count int
	var b_name string
	test := make(map[string]interface{})
	// test["delete"] = make([]string, 0)

	branch_count := "SELECT max(branch_num) from basic where name=1"
	row, err := database.DB().Query(branch_count)
	utils.HandleError(err)
	defer row.Close()

	for row.Next() {
		row.Scan(&b_count)
	}
	// fmt.Println(b_count)

	for i := 1; i <= b_count; i++ {
		// print(i)
		branch_num := fmt.Sprintf("Branch%02d", i)
		query := "SELECT branch, payload, abilityname from basic where name=1 AND branch=? ORDER BY branch, seq"
		// print(query)
		rows, err := database.DB().Query(query, branch_num)
		// print(rows)

		utils.HandleError(err)
		defer rows.Close()

		var data []Send

		for rows.Next() {

			var pay string
			var abname string
			rows.Scan(&b_name, &pay, &abname)
			// fmt.Println(b_name)
			data = append(data, Send{pay, abname})
			// fmt.Println(data)

		}
		// fmt.Println(b_name)
		test[b_name] = data
	}
	fmt.Println(("잘 되는 중~"))
	dataBytes, _ := json.Marshal(test)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataBytes)
	// json.NewEncoder(w).Encode(dataBytes)
	// fmt.Println(string(dataBytes))

}
