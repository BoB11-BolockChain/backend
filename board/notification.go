package board

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	_ "github.com/mattn/go-sqlite3"
)

type CrudNotify struct {
	Num  int    `json:"num"`
	Crud string `json:"crud"`
}

type Notifications struct {
	Num         int    `json:"num"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Author      string `json:"author"`
	CreatedDate string `json:"cdate"`
	Views       int    `json:"views"`
}

func Notification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	if r.Method == "POST" {
		var crudNoti CrudNotify
		json.NewDecoder(r.Body).Decode(&crudNoti)
		fmt.Println(crudNoti.Crud)
		switch crudNoti.Crud {
		case "Edit":
			return
		case "Remove":
			query := "delete from notification where num=?"
			fmt.Println(query)
			_, err := database.DB().Exec(query, crudNoti.Num)
			if err != nil {
				panic(err)
			}
			return
		default:
			return
		}
	}
	// fmt.Fprint(w, r.Form)
	var noti_ret Notifications
	var noti_retn []Notifications
	query := "select num, title, content, author, cdate, views from notification order by num desc limit ?,1"

	row, err := database.DB().Query("select count(*) from notification")
	if err != nil {
		panic(err)
	}
	defer row.Close()
	var count int

	for row.Next() {
		row.Scan(&count)
	}
	for i := 0; i < count; i++ {

		fmt.Println(i)
		row := database.DB().QueryRow(query, i)
		switch err := row.Scan(&noti_ret.Num, &noti_ret.Title, &noti_ret.Content, &noti_ret.Author, &noti_ret.CreatedDate, &noti_ret.Views); err {
		case nil:
			noti_retn = append(noti_retn, noti_ret)
		default:
			panic(err)
		}

	}
	data := struct {
		Data []Notifications `json:"data"`
	}{noti_retn}
	json.NewEncoder(w).Encode(data)
	fmt.Println(data)

	fmt.Println("게시글 정보 전송 완료")
}
