package challenges

import (
	"fmt"
	"net/http"
	// "database/sql"
	// "encoding/json"


	"github.com/backend/database"
	"github.com/backend/utils"
	"github.com/blockloop/scan"
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
	Num int			`json:"num"`
	Title string	`json:"title"`
	Score string	`json:"score"`
}

func ViewInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var info View
	query := "select num, title, score from challenges where num=?"
	rows, err := database.DB().Query("select count(*) from challenges")
	utils.HandleError(err)
	defer rows.Close()
	var count int

	for rows.Next(){		
		rows.Scan(&count)
	}

	for i :=1; i <= count; i++ {
		fmt.Println(i)
		row := database.DB().QueryRow(query, i)
		switch err := row.Scan(&info.Num, &info.Title, &info.Score); err{
		case nil:
		fmt.Println(info)
		// &info.Num, &info.Title, &info.Score
		// enc := json.NewEncoder(w)
		// w.Header().Set("Content-Type", "application/json")
		// enc.Encode(u)
		default:
		panic(err)
	}

	
	// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
	return

	}
	

	// for value=1; value <= 
	// info, err := database.DB().Query("select num, title, score from challenges")where num=?", value)

}


func ChInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var challenges Challenge
	var num int
	value := r.FormValue("num")
	// json.Unmarshal([]byte(), &num)
	fmt.Println(value)
	fmt.Println(num)
	ch, err := database.DB().Query("select * from challenges where num=?", value)
	utils.HandleError(err)
	fmt.Println(ch)

	ch_err := scan.Row(&challenges, ch)
	fmt.Println(ch_err)
	if ch_err != nil {
		// return challenges{}, err
	}
	// return challenges, nil

}
