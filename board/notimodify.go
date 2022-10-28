package board

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/backend/database"
	_ "github.com/mattn/go-sqlite3"
)

type Noticreate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"id"`
}

type Notiedit struct {
	Num     int    `json:"num"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"id"`
}

func NotiCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	if r.Method == "POST" {
		var createNoti Noticreate
		json.NewDecoder(r.Body).Decode(&createNoti)
		fmt.Println(createNoti)
		now := time.Now()
		custom_now := now.Format("2006-01-02 15:04:05")
		fmt.Println(custom_now)
		insert, _ := database.DB().Prepare("INSERT INTO notification (title, content, author, cdate, views) VALUES (?, ?, ?, ?, 0)")
		_, err := insert.Exec(createNoti.Title, createNoti.Content, createNoti.Author, custom_now)
		if err != nil {
			panic(err)
			//utils.HandleError(err)
		}
	}
}

func NotiEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	if r.Method == "POST" {
		var editNoti Notiedit
		json.NewDecoder(r.Body).Decode(&editNoti)
		fmt.Println(editNoti)
		now := time.Now()
		custom_now := now.Format("2006-01-02 15:04:05")
		fmt.Println(custom_now)
		insert, _ := database.DB().Prepare("UPDATE notification SET title = ?, content = ? WHERE num = ?")
		_, err := insert.Exec(editNoti.Title, editNoti.Content, editNoti.Num)
		if err != nil {
			panic(err)
			//utils.HandleError(err)
		}
	}
}
