package challenges

import (
	"fmt"
	"database/sql"

	"github.com/backend/utils"
	"github.com/backend/database"
)

type Challenge struct {
	Num int
	Title string
	Desc string
	Os string
	Score string
	Attack string
}

func challenge () {
	ch, err := database.DB().Query("select * from challenges")
	HandleError(err)
}