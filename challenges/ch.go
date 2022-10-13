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
	Num    int
	Title  string
	Desc   string
	Os     string
	Score  string
	Attack string
}

type View struct {
	Num   int    `json:"num"`
	Title string `json:"title"`
	Score string `json:"score"`
}

type Number struct {
	Num int `json:"num"`
}

func ViewInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var info View
	var sendInfo []View
	query := "select num, title, score from challenges where num=?"
	rows, err := database.DB().Query("select count(*) from challenges")
	utils.HandleError(err)
	defer rows.Close()
	var count int

	for rows.Next() {
		rows.Scan(&count)
	}

	fmt.Println(count)

	for i := 1; i <= count; i++ {

		fmt.Println(i)
		row := database.DB().QueryRow(query, i)
		switch err := row.Scan(&info.Num, &info.Title, &info.Score); err {
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

	query := "select num, title, desc, os, score, attack from challenges where num=?"
	print(query)
	row := database.DB().QueryRow(query, value)
	print(row)
	err = row.Scan(&challenges.Num, &challenges.Title, &challenges.Desc, &challenges.Attack, &challenges.Os, &challenges.Score)
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
