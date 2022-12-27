package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/backend/database"
	"github.com/backend/utils"
	"github.com/google/uuid"
)

var Sessions = map[string]session{}

type session struct {
	userid   string
	expiry   time.Time
	auth     string
	vmname   string
	vmexpiry time.Time
}

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

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
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
	}

	if email_err == sql.ErrNoRows && id_err == sql.ErrNoRows {
		insert, _ := database.DB().Prepare("INSERT INTO user (email, id, pw, auth) values(?, ?, ?, 'user')")
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
	// fmt.Fprint(w, r.Form)
	var new Login
	json.NewDecoder(r.Body).Decode(&new)
	fmt.Println(new)
	pwhash := utils.Hash(new.Pw)

	query := fmt.Sprintf("SELECT COUNT(*) as count, auth FROM user where id='%s' and pw='%s'", new.Id, pwhash)
	rows, err := database.DB().Query(query)
	fmt.Println(query)
	utils.HandleError(err)
	defer rows.Close()

	enc := json.NewEncoder(w)

	for rows.Next() {
		var count int
		var auth string
		rows.Scan(&count, &auth)

		if count == 1 {
			//success
			// Create a new random session token
			sessionToken := uuid.NewString()
			expiresAt := time.Now().Add(180 * time.Second)

			Sessions[sessionToken] = session{
				userid:   new.Id,
				expiry:   expiresAt,
				auth:     auth,
				vmname:   "",
				vmexpiry: expiresAt,
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Path:    "/admin",
				Expires: expiresAt,
			})

			fmt.Println(sessionToken)

			userSession, exists := Sessions[sessionToken]
			if !exists {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if userSession.isExpired() {
				delete(Sessions, sessionToken)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			fmt.Println(userSession)
			data := struct {
				SessionId string `json:"sessionId"`
			}{userSession.userid}

			enc.Encode(data)
			// fmt.Println(data)
			fmt.Println("success!")
		} else {
			//fail
			// fmt.Println("login fail")
			fmt.Fprint(w, "로그인에 실패했습니다")
			type Fail struct {
				Status string `json:"status"`
			}
			u := Fail{Status: "fail"}
			w.Header().Set("Content-Type", "application/json")
			// enc.Encode(u)
			json.NewEncoder(w).Encode(u)
			// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
			// return
		}
	}
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value
	userSession, exists := Sessions[sessionToken]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(Sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("welcome, %s!", userSession.userid)))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	delete(Sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_cookie",
		Value:   "",
		Expires: time.Now(),
	})
}
