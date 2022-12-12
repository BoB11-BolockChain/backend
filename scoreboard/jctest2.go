package scoreboard

import (
	"fmt"

	"github.com/backend/database"
	_ "github.com/mattn/go-sqlite3"
)

func Cal() {
	//aaa := caldera.GetOperation("ab5647d5-09d1-4389-a30f-5fad9f5dee94") // id 입력
	var a string
	fmt.Println("1")
	query := "select user_id from solved_challenge where user_id = 'wrwrwrwr' and solved_challenge_id = 40"
	row := database.DB().QueryRow(query)
	row.Scan(&a)
	fmt.Println(a)
}
