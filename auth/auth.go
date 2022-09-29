package auth

import (
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

var sessions map[string]string

func SignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprint(w, r.Form)
	formData := r.PostForm
	id := formData.Get("id")
	pw := formData.Get("pw")
	pwhash := utils.Hash(pw)

	query := fmt.Sprintf("SELECT COUNT(*) as count FROM user where id='%s' and pw='%s'", id, pwhash)
	rows, err := database.DB().Query(query)
	utils.HandleError(err)
	defer rows.Close()

	for rows.Next() {
		var count int
		rows.Scan(&count)

		if count == 1 {
			//success
			sessions[id] = id
			http.Redirect(w, r, "URL_TO_MAIN_PAGE", http.StatusSeeOther)
		} else {
			//fail
			fmt.Fprint(w, "hello go")
			http.Redirect(w, r, "URL_TO_LOGIN_PAGE", http.StatusSeeOther)
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
