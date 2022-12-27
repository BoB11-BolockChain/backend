package challenges

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/backend/database"
	"github.com/backend/utils"
)

type Ability struct {
	Data []struct {
		Branch      string `json:"branch"`
		Seq         string `json:"seq"`
		Payload     string
		AbilityName string
	} `json:"data"`
}

type ChallengeNum struct {
	Num string
}

type SendAbility struct {
	Payload     string
	AbilityName string
}

type SendBranch struct {
	Data []struct {
		Payload     string
		AbilityName string
	}
}

func GetChNum(w http.ResponseWriter, r *http.Request) {

}



func InsertData2(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	var ability Ability
	var _chnum int

	json.NewDecoder(r.Body).Decode(&ability)

	insert, _ := database.DB().Prepare("INSERT INTO ability (ch_num, branch_num, seq, abilityname, payload) values(?, ?, ?, ?, ?)")

	count := len(ability.Data)

	for i := 0; i < count; i++ {
		_branch := ability.Data[i].Branch
		_seq := ability.Data[i].Seq
		_abilityname := ability.Data[i].AbilityName
		_payload := ability.Data[i].Payload
		_, err := insert.Exec(_chnum, _branch, _seq, _abilityname, _payload)
		if err != nil {
			utils.HandleError(err)
		}
		fmt.Println(err)
	}
}

func PrintData(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var ch_num ChallengeNum
	var data SendAbility
	json.NewDecoder(r.Body).Decode(&ch_num)
	fmt.Println(ch_num)

	// get branch_num (1, 2, 3)
	query := fmt.Sprintf("select payload, abilityname from branch  where ch_num=%s ORDER BY branch_num, seq", ch_num)
	print(query)
	rows, err := database.DB().Query(query)
	fmt.Println(err)
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&data.Payload, &data.AbilityName)
		fmt.Println(data)
	}

	// query := fmt.Sprintf("select payload, abilityname from ability where ='%s'", session.SessionId)

}
