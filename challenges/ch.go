package challenges

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	// "reflect"

	// "database/sql"
	"encoding/json"

	"github.com/backend/database"
	"github.com/backend/utils"

	// "github.com/blockloop/scan"
	// "github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Challenge struct {
	Num   int
	Title string
	Desc  string
	Os    string
	Score string
}

type View struct {
	Num   int    `json:"num"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
	Score int    `json:"score"`
	Os    string `json:"os"`
}

type Number struct {
	Num int `json:"num"`
}

type DeleteChall struct {
	Num  int    `json:"num"`
	Crud string `json:"crud"`
}

func ViewInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	if r.Method == "POST" {
		var deletechall DeleteChall
		json.NewDecoder(r.Body).Decode(&deletechall)
		if deletechall.Crud == "Remove" {
			// query := "delete from challenge where num=?"
			// _, err := database.DB().Exec(query, deletechall.Num)

			branch_query := "select num from branch where chall_num=?"
			row, err := database.DB().Query(branch_query, deletechall.Num)
			if err != nil {
				panic(err)
			}
			defer row.Close()

			var count int
			for row.Next() {
				row.Scan(&count)
				fmt.Println(count)
				fmt.Println("--")
			}

			// query := "delete from branch where chall_num=?"
			// _, err := database.DB().Exec(query, deletechall.Num)
			// query := "delete from ability where branch_num=?"
			// _, err := database.DB().Exec(query, deletechall.Num)
		}
	}

	var info View
	var sendInfo []View
	query := "select num, title, desc, score, os from challenge where num=?"
	rows, err := database.DB().Query("select count(*) from challenge")
	utils.HandleError(err)
	defer rows.Close()
	var count int
	// rows.Scan(&count)
	// fmt.Println(count)
	for rows.Next() {
		rows.Scan(&count)
	}

	fmt.Println(count)

	for i := 1; i <= count; i++ {

		fmt.Println(i)
		row := database.DB().QueryRow(query, i)
		switch err := row.Scan(&info.Num, &info.Title, &info.Desc, &info.Score, &info.Os); err {
		case nil:
			sendInfo = append(sendInfo, info)
		default:
			panic(err)
		}

	}
	b, err := json.Marshal(sendInfo)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	fmt.Println(string(b))
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Data []View `json:"data"`
	}{sendInfo}
	json.NewEncoder(w).Encode(data)
	fmt.Println(data)

	// json.NewEncoder(w).Encode(string(b))
}

func ChInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var challenges Challenge
	var value int

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	utils.HandleError(err)
	value = id
	fmt.Println(value)

	query := "select num, title, desc, os, score from challenges where num=?"
	print(query)
	row := database.DB().QueryRow(query, value)
	print(row)
	err = row.Scan(&challenges.Num, &challenges.Title, &challenges.Desc, &challenges.Os, &challenges.Score)
	print(err)

	b, err := json.Marshal(challenges)
	if err != nil {
		log.Fatalf("JSON marshaling failed: %s", err)
	}

	fmt.Println(string(b))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(string(b))

	// json.NewEncoder(w).Encode(challenges)
}
