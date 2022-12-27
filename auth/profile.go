package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/backend/database"
	"github.com/backend/makevm"
	"github.com/backend/utils"
)

type UserProfile struct {
	Id      string
	Email   string
	Org     sql.NullString
	Comment sql.NullString
}

type RE_UserProfile struct {
	Id      string
	Email   string
	Org     string
	Comment string
}

type Id struct {
	Id string
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	// var userId string
	var userInfo UserProfile
	var userId Id
	json.NewDecoder(r.Body).Decode(&userId)
	// fmt.Println(userId.Id)

	userInfo.Id = userId.Id

	select_query := "select email, organization, comment FROM user WHERE id=?"
	err := database.DB().QueryRow(select_query, userInfo.Id).Scan(&userInfo.Email, &userInfo.Org, &userInfo.Comment)
	utils.HandleError(err)

	fmt.Println(userInfo)
	json.NewEncoder(w).Encode(userInfo)
}

func SaveProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	var userInfo RE_UserProfile
	json.NewDecoder(r.Body).Decode(&userInfo)
	fmt.Print("수신:")
	fmt.Println(userInfo)

	insert, _ := database.DB().Prepare("update user set organization=?, comment=? where id=?")
	_, err := insert.Exec(userInfo.Org, userInfo.Comment, userInfo.Id)
	if err != nil {
		panic(err)
	}

	fmt.Println("완료~")

	time.Sleep(1000 * time.Millisecond)
	dirname := "/home/ar/user_windows/profile/"
	makevm.ExcuteCMD("sudo", "sh", "-c", "mv -f "+dirname+"temp"+" /var/myapp/public/Profile/"+userInfo.Id+".png")
}

func UploadImg(w http.ResponseWriter, r *http.Request) {
	uploadFile, header, err := r.FormFile("upfiles")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	defer uploadFile.Close()
	if header == nil {
		fmt.Println("no Header")
	}
	dirname := "/home/ar/user_windows/profile"
	os.MkdirAll(dirname, 0777)
	filepath := fmt.Sprintf("%s/%s", dirname, "temp")
	file, err := os.Create(filepath)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	io.Copy(file, uploadFile)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filepath)
	defer file.Close()
}
