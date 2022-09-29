package challenge

import (
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

func GetChallengeById(w http.ResponseWriter, r *http.Request) {
	var values string
	err := database.DB().QueryRow("select * from challenges where id=$1", r.URL.Query().Get("id")).Scan(&values)
	utils.HandleError(err)

	fmt.Fprintf(w, "%s", values)
}
