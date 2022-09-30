package auth

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/backend/database"
	"github.com/backend/utils"
)

var sessions map[string]string

func SignUp(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fmt.Fprint(w, r.Form)
	email := r.FormValue("email")
	id := r.FormValue("id")
	pw := r.FormValue("pw")
	pwhash := utils.Hash(pw)

	fmt.Println(email)
	fmt.Println(id)
	fmt.Println(pw)
	fmt.Println(pwhash)

	var user string

	email_err := database.DB().QueryRow("SELECT email From user WHERE email=?", email).Scan(&user)
	id_err := database.DB().QueryRow("SELECT id From user WHERE id=?", id).Scan(&user)

	if email_err == sql.ErrNoRows && id_err == sql.ErrNoRows {
		insert, _ := database.DB().Prepare("INSERT INTO user (email, id, pw) values(?, ?, ?)")
		_, err := insert.Exec(email, id, pwhash)
		if err != nil {
			utils.HandleError(email_err)
		}
		fmt.Fprint(w, "hi!")
		fmt.Println("성공!")
		return
	// } else if email_err != nil {
	// 	utils.HandleError(email_err)
	// 	fmt.Println("존재하는 이메일입니다")
	// 	fmt.Fprint(w, "존재하는 이메일입니다")
	// 	// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
	// } else if id_err != nil {
	// 	utils.HandleError(id_err)
	// 	fmt.Println("존재하는 아이디입니다")
	// 	fmt.Fprint(w, "존재하는 아이디입니다")
	// 	// http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
	} else {
		if email_err != nil {
			fmt.Println("1")
		}
		if id_err != nil {
			fmt.Println("2")
		}
		fmt.Println("회원가입 실패")
		utils.HandleError(email_err)
		utils.HandleError(id_err)
	}
}


func SignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprint(w, r.Form)
	id := r.FormValue("id")
	pw := r.FormValue("pw")
	pwhash := utils.Hash(pw)

	query := fmt.Sprintf("SELECT COUNT(*) as count FROM user where id='%s' and pw='%s'", id, pwhash)
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

// if($count>0){
//     if(isset($id)){
//         $user_login = TRUE;
//         $_SESSION['id'] = $id;
//         echo "<script>alert('login as $id')</script>";
//         echo "<meta http-equiv='refresh' content='0;url=../admin.index.php'>";
//     }else{
//         echo "login fail!";
//         echo "<meta http-equiv='refresh' content='0;url=./signin.php'>";
//     }
// }else if($num == 0){
//     echo "<script>alert('No information!');</script>";
//     echo "<meta http-equiv='refresh' content='0;url=./signup.php'>";
// }else{
//     echo "<script>alert('Error');</script>";
//     echo "<meta http-equiv='refresh' content='0;url=./signin.php'>";
// }
