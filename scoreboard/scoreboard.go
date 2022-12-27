package scoreboard

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"sort"

	"github.com/backend/database"
	// "github.com/backend/utils"
	_ "github.com/mattn/go-sqlite3"
)

type ScoreValue struct {
	Num     int        `json:"num"`
	Id      string     `json:"id"`
	Score   int        `json:"score"`
	Total int `json:"total"`
	IR      int        `json:"IR"`
	Org     string     `json:"org"`
	Comment string     `json:"comment"`
	Solved  []SolvedCh `json:"solved"`
}

type IrScore struct {
	UserId string `json:"userid"`
	Sce_Id int    `json:"sce_id"`
	Score  int    `json:"score"`
}

type SendAll struct {
	Data []ScoreValue `json:"data"`
	Line []GraphValue `json:"line"`
}

type UserId struct {
	UserId string `json:"userId"`
}

type SolvedCh struct {
	Scenario_title string      `json:"scenario_title"`
	IrScore        int         `json:"irscore"`
	Challenges     []Challenge `json:"challenge"`
}

type Challenge struct {
	Challenge_title string `json:"challenge_title"`
	Time            string `json:"time"`
	Score           int    `json:"score"`
}

type ModalValue struct {
	Num             int    `json:"num"`
	Scenario_id     int    `json:"scenario_id"`
	Scenario_title  string `json:"scenario_title"`
	Challenge_title string `json:"challenge_title"`
	Time            string `json:"time"`
	Score           int    `json:"score"`
}

type GraphValue struct {
	Time  string `json:"time"`
	Score int    `json:"score"`
	User  string `json:"user"`
}

func GetScore(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	// fmt.Println("전송 시작!")

	// var allData SendAll
	var scoredata ScoreValue
	var scoresend []ScoreValue
	var count int
	var graphdata GraphValue
	var graphsend []GraphValue

	count_query := "SELECT COUNT(DISTINCT user_id) FROM solved_challenge"
	row := database.DB().QueryRow(count_query)
	row.Scan(&count)
	// fmt.Println(count)
	query := "SELECT sum(score), user_id from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id group by user_id ORDER BY sum(score) DESC limit ?,1;"
	sce_count_query := "SELECT count(scenario_id) from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=?"
	ch_count_query := "SELECT count(*) from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? and scenario_id=?"
	sce_query := "SELECT DISTINCT scenario_id from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? limit ?,1"
	ch_query := "SELECT title, solved_time, score from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? and scenario_id=? ORDER by solved_time limit ?,1"
	sce_title_query := "SELECT title from scenario where id=?"

	top_query := "SELECT user_id from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id group by user_id ORDER BY sum(score) DESC limit ?, 1;"
	line_query := "SELECT score, solved_time from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? ORDER by solved_time limit ?,1;"
	point_query := "SELECT count(*) from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? "

	additional_info_query := "SELECT organization, comment from user where id=?"

	ir_score_query := "SELECT score from solved_ir_score where user_id=? and sce_id=?"

	// sum_query := "SELECT sum(challenge.score), sum(solved_ir_score), user_id from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id group by user_id ORDER BY sum(score) DESC limit ?,1;"

	inital_time_query := "SELECT solved_time from solved_challenge ORDER BY solved_time limit 1"
	last_time_query := "SELECT solved_time from solved_challenge ORDER BY solved_time DESC limit 1"
	var initial_time string
	var last_time string
	var ir_score_sum int
	row = database.DB().QueryRow(inital_time_query)
	inverse_row := database.DB().QueryRow(last_time_query)
	row.Scan(&initial_time)
	inverse_row.Scan(&last_time)
	// fmt.Println(initial_time)
	// fmt.Println(last_time)
	time_parsing, err := time.Parse("2006-01-02 15:04:05", initial_time)
	if err != nil {
		panic(err)
	}
	convMinutes, _ := time.ParseDuration("10m")
	time_deduct := time_parsing.Add(-convMinutes).Format("2006-01-02 15:04:05")
	// fmt.Println(time_deduct)

	if 10 < count {
		count = 10
	}

	for i := 0; i < count; i++ {
		scoredata.Num = i + 1
		row := database.DB().QueryRow(query, i)
		
		switch err := row.Scan(&scoredata.Score, &scoredata.Id); err {
		case nil:
			// fmt.Println(scoredata)
			row := database.DB().QueryRow(sce_count_query, scoredata.Id)
			row2 := database.DB().QueryRow(additional_info_query, scoredata.Id)
			row2.Scan(&scoredata.Org, &scoredata.Comment)
			var sce_count int
			row.Scan(&sce_count)
			// fmt.Println("여기에 정보 >> ")
			for j := 0; j < sce_count; j++ {
				var sce_id int
				row := database.DB().QueryRow(sce_query, scoredata.Id, j)

				switch err := row.Scan(&sce_id); err {
				case nil:

					var chdatapack []Challenge
					var solvedpack SolvedCh
					var ch_count int
					var chdata Challenge
					title_row := database.DB().QueryRow(sce_title_query, sce_id)
					row2 := database.DB().QueryRow(ir_score_query, scoredata.Id, sce_id)
					row2.Scan(&solvedpack.IrScore)
					// fmt.Println("여 기 나 와")
					ir_score_sum += solvedpack.IrScore
					// fmt.Println(solvedpack.IrScore)

					title_row.Scan(&solvedpack.Scenario_title)
					ch_count_row := database.DB().QueryRow(ch_count_query, scoredata.Id, sce_id)
					ch_count_row.Scan(&ch_count)
					for k := 0; k < ch_count; k++ {
						ch_row := database.DB().QueryRow(ch_query, scoredata.Id, sce_id, k)
						ch_row.Scan(&chdata.Challenge_title, &chdata.Time, &chdata.Score)
						// fmt.Println(chdata)
						chdatapack = append(chdatapack, chdata)
						// fmt.Println(chdatapack)
					}
					solvedpack.Challenges = chdatapack
					// fmt.Println(solvedpack)
					scoredata.Solved = append(scoredata.Solved, solvedpack)
					// fmt.Println(scoredata.Solved)
					scoredata.IR = ir_score_sum
					ir_score_sum = 0
				}
			}
			scoredata.Total = scoredata.Score + scoredata.IR
			scoresend = append(scoresend, scoredata)
			scoredata = ScoreValue{}
		case sql.ErrNoRows:
			//park
			continue
		default:
			panic(err)
		}
	}

	for i := 0; i < count; i++ {
		var top_id string
		row := database.DB().QueryRow(top_query, i)
		switch err := row.Scan(&top_id); err {
		case nil:
			point_row, err := database.DB().Query(point_query, top_id)
			if err != nil {
				panic(err)
			}
			graphdata.Time = time_deduct
			graphdata.User = top_id
			graphdata.Score = 0

			graphsend = append(graphsend, graphdata)
			for point_row.Next() {
				var point_count int
				var score int
				var score_sum int
				point_row.Scan(&point_count)
				for j := 0; j < point_count; j++ {
					rows := database.DB().QueryRow(line_query, top_id, j)
					switch err := rows.Scan(&score, &graphdata.Time); err {
					case nil:
						score_sum = score_sum + score
						graphdata.Score = score_sum
						graphdata.User = top_id
						// fmt.Println(graphdata)
						graphsend = append(graphsend, graphdata)
					default:
						panic(err)
					}
				}
				graphdata.Score = score_sum
				graphdata.Time = time.Now().Format("2006-01-02 15:04:05")
				graphsend = append(graphsend, graphdata)
				score_sum = 0
			}
		}
	}

	sort.Slice(scoresend, func(i, j int) bool {
		return scoresend[i].Total > scoresend[j].Total
	})

	// fmt.Println(graphsend)
	data := struct {
		Data []ScoreValue `json:"data"`
		Line []GraphValue `json:"line"`
	}{scoresend, graphsend}

	json.NewEncoder(w).Encode(data)
}

func GetScoreModal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	fmt.Println("개인 데이터 전송 시작!")
	var getId UserId
	var userId string
	var num int
	var senduserdata []ModalValue

	fmt.Println("checkpoint")
	json.NewDecoder(r.Body).Decode(&getId)
	fmt.Println(getId.UserId)

	userId = string(getId.UserId)
	fmt.Println(userId)

	pre_query := "SELECT count(*) from solved_challenge where user_id=?"
	row, err := database.DB().Query(pre_query, userId)
	if err != nil {
		panic(err)
	}
	defer row.Close()

	for row.Next() {
		row.Scan(&num)
		println(num)
	}

	query := "SELECT title, score, scenario_id, solved_time from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? ORDER by solved_time limit ?,1"

	second_query := "SELECT title from scenario where id=?"

	for i := 0; i < num; i++ {
		var userdata ModalValue
		// userId = "aaa"
		println(i)
		userdata.Num = i + 1
		row := database.DB().QueryRow(query, userId, i)

		switch err := row.Scan(&userdata.Challenge_title, &userdata.Score, &userdata.Scenario_id, &userdata.Time); err {
		case nil:
			fmt.Println(userdata)
			second_row := database.DB().QueryRow(second_query, userdata.Scenario_id)
			err := second_row.Scan(&userdata.Scenario_title)
			if err != nil {
				panic(err)
			}
			senduserdata = append(senduserdata, userdata)
		default:
			panic(err)
		}
	}

	fmt.Println(senduserdata)
	data := struct {
		Data []ModalValue `json:"data"`
	}{senduserdata}
	json.NewEncoder(w).Encode(data)
	// fmt.Println("모달 데이터 가 는 중 ~ ")
}

func GetGraphData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	fmt.Println("전송 시작!")

	var graphdata GraphValue
	var graphsend []GraphValue
	var count int
	count_query := "SELECT COUNT(DISTINCT user_id) FROM solved_challenge"
	row := database.DB().QueryRow(count_query)
	row.Scan(&count)
	fmt.Println(count)

	top_query := "SELECT user_id from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id group by user_id ORDER BY sum(score) DESC limit ?, 1;"

	query := "SELECT score, solved_time from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? ORDER by solved_time limit ?,1;"

	// point_count
	point_query := "SELECT count(*) from solved_challenge inner join challenge on solved_challenge.solved_challenge_id = challenge.id where user_id=? "

	if 10 < count {
		count = 10
	}

	for i := 0; i < count; i++ {
		var top_id string
		row := database.DB().QueryRow(top_query, i)
		switch err := row.Scan(&top_id); err {
		case nil:
			fmt.Println(top_id)
			point_row, err := database.DB().Query(point_query, top_id)
			if err != nil {
				panic(err)
			}
			for point_row.Next() {
				var point_count int
				point_row.Scan(&point_count)
				fmt.Println(point_count)
				for j := 0; j < point_count; j++ {
					rows := database.DB().QueryRow(query, top_id, j)
					switch err := rows.Scan(&graphdata.Score, &graphdata.Time); err {
					case nil:
						graphdata.User = top_id
						fmt.Println(graphdata)
						graphsend = append(graphsend, graphdata)
					default:
						panic(err)
					}
				}
			}
		}
	}
	fmt.Println(graphsend)
	line := struct {
		Data []GraphValue `json:"line"`
	}{graphsend}
	json.NewEncoder(w).Encode(line)
	fmt.Println("check check")
}
