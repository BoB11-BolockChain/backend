package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	// "github.com/backend/utils"
	// "github.com/backend/auth/auth"
)

type Profile struct {
	UserId    string `json:"id"`
	UserEmail string `json:"email"`
}

type Ses struct {
	SessionId string `json:"sessionid"`
}

func UserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var session Ses
	json.NewDecoder(r.Body).Decode(&session)
	fmt.Println(session)
	var profile Profile
	query := fmt.Sprintf("select id, email from user where id='%s'", session.SessionId)
	print(query)
	row := database.DB().QueryRow(query)
	switch err := row.Scan(&profile.UserId, &profile.UserEmail); err {
	case nil:
		w.Header().Set("Content-Type", "application/json")
		data := struct {
			Data Profile `json:"data"`
		}{profile}
		json.NewEncoder(w).Encode(data)
	default:
		panic(err)
	}
}
