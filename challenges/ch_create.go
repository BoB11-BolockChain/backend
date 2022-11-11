package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
	_ "github.com/mattn/go-sqlite3"
)

// type data struct {
// 	Data []struct {
// 		Branch []struct {
// 			Seq         int    `json:"Sequence"`
// 			Payload     string `json:"Payload"`
// 			Abilityname string `json:"AbilityName"`
// 		} `json:"Branch"`
// 		//Delete string `json:Delete`
// 	} `json:"data"`
// 	Info struct {
// 		Title       string `json:"Title"`
// 		Description string `json:"Description"`
// 		Scenario_id       int    `json:"Scenario_id"`
// 		Sequence          string `json:"Sequence"`
// 	} `json:"info"`
// }

type Scenario struct {
	Title     string `json:"title"`
	Desc      string `json:"desc"`
	System        string `json:"system"`
	VMOptions []struct {
	   Name    string `json:"name"`
	   Command string `json:"command"`
	} `json:"vm-options"`
	Data []struct {
	   Tactic   string   `json:"tactic"`
	   Payloads []string `json:"payloads"`
	} `json:"data"`
 }

// type Send struct {
// 	Payload     string `json:"Payload"`
// 	Abilityname string `json:"AbilityName"`
// }

func InsertData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	if r.Method == "POST" {
		var ch_dat data

		json.NewDecoder(r.Body).Decode(&ch_dat)
		fmt.Println(ch_dat)
		fmt.Println(ch_dat.Data)
		sce_insert, _ := database.DB().Prepare("INSERT INTO scenario (id, title, description, system) VALUES (NULL, ?, ?, ?)")
		_, err := sce_insert.Exec(ch_dat.Scenario.Title, ch_dat.Scenario.Description, ch_dat.Scenario.System)
		if err != nil {
			utils.HandleError(err)
		}

		vm_insert, _ := database.DB().Prepare("INSERT INTO vm_option (id, title, command) VALUES (NULL, ?, ?)")
		_, err := vm_insert.Exec(ch_dat.Scenario.Data.Name, ch_dat.Scenario.Data.Command)
		if err != nil {
			utils.HandleError(err)
		}
		
		chall_insert, _ := database.DB().Prepare("INSERT INTO challenge (id, title, description, scenario_id, sequence) VALUES (NULL, ?, ?, ?, ?)")
		_, err := chall_insert.Exec(ch_dat.Scenario.Title, ch_dat.Scenario.Description, ch_dat.Scenario.Scenario_id, ch_dat.Scenario.Sequence)
		if err != nil {
			utils.HandleError(err)
		}

		pay_insert, _ := database.DB().Prepare("INSERT INTO challenge (id, title, description, scenario_id, sequence) VALUES (NULL, ?, ?, ?, ?)")
		_, err := chall_insert.Exec(ch_dat.Scenario.Title, ch_dat.Scenario.Description, ch_dat.Scenario.Scenario_id, ch_dat.Scenario.Sequence)
		if err != nil {
			utils.HandleError(err)
		}

		var sce_num int
		sce_num_query := "SELECT id FROM scenario WHERE title=? ORDER BY rowid DESC LIMIT 1"
		row := database.DB().QueryRow(sce_num_query, ch_dat.makeScenario.Title)
		row.Scan(&sce_num)

		for i, challenge := range ch_dat.Data {
			challenge_insert, _ := database.DB().Prepare("INSERT INTO challenge (id, title, description, scenario_id, sequence) VALUES (NULL, ?, ?, ?)")
			_, err := challenge_insert.Exec(chall_num, i+1)
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
	// fmt.Println(("잘 되는 중~"))
	dataBytes, _ := json.Marshal(test)
	w.Header().Set("Content-Type", "application/json")
	w.Write(dataBytes)
	// json.NewEncoder(w).Encode(dataBytes)
	// fmt.Println(string(dataBytes))

}
