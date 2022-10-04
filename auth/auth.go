package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

var sessions map[string]string

type User struct {
	Email string
	Id    string
	Pw    string
	Conpw string
}

type Login struct {
	Id string
	Pw string
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	// fmt.Fprint(w, r.Form)
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)

	pwhash := utils.Hash(user.Pw)

	var newuser string

	email_err := database.DB().QueryRow("SELECT email From user WHERE email=?", user.Email).Scan(&newuser)
	id_err := database.DB().QueryRow("SELECT id From user WHERE id=?", user.Id).Scan(&newuser)

	if user.Conpw != user.Pw {
		fmt.Fprint(w, "비밀번호가 일치하지 않습니다.")
		// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)

	}

	if email_err == sql.ErrNoRows && id_err == sql.ErrNoRows {
		insert, _ := database.DB().Prepare("INSERT INTO user (email, id, pw) values(?, ?, ?)")
		_, err := insert.Exec(user.Email, user.Id, pwhash)
		if err != nil {
			utils.HandleError(email_err)
		}
		fmt.Fprint(w, "hi!")
		fmt.Println("성공!")
		return
	} else {
		if email_err == nil {
			fmt.Println("이미 존재하는 이메일입니다")
			// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
		}
		if id_err == nil {
			fmt.Println("이미 존재하는 아이디입니다")
			// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
		}
		fmt.Println("회원가입 실패")
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	fmt.Fprint(w, r.Form)
	var new Login
	json.NewDecoder(r.Body).Decode(&new)
	fmt.Println(new)
	pwhash := utils.Hash(new.Pw)

	query := fmt.Sprintf("SELECT COUNT(*) as count FROM user where id='%s' and pw='%s'", new.Id, pwhash)
	rows, err := database.DB().Query(query)
	fmt.Println(query)
	utils.HandleError(err)
	defer rows.Close()

	for rows.Next() {
		var count int
		rows.Scan(&count)

		if count == 1 {
			//success
			// sessions[id] = id
			fmt.Println("hihi")
			// http.Redirect(w, r, "URL_TO_MAIN_PAGE", http.StatusSeeOther)
		} else {
			//fail
			fmt.Println("login fail")
			fmt.Fprint(w, "로그인에 실패했습니다")
			// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
		}
	}
}
