package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/backend/auth"
	"github.com/backend/board"
	"github.com/backend/caldera"
	"github.com/backend/create"
	"github.com/backend/dashboard"

	"github.com/backend/makevm"
	"github.com/backend/scoreboard"
	"github.com/backend/training"
	"github.com/backend/utils"
	"github.com/gorilla/mux"
)

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		next.ServeHTTP(w, r)
	})
}

func optionMethodBanMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			next.ServeHTTP(w, r)
		} else {
			w.Write([]byte("message from option block middleware"))
		}
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Print("hello")
}

func getabs(w http.ResponseWriter, r *http.Request) {
	addr := "http://www.pdxf.tk:8888/api/v2/abilities"
	req, err := http.NewRequest("GET", addr, nil)
	utils.HandleError(err)

	req.Header.Add("KEY", "ADMIN123")

	client := http.Client{}
	res, err := client.Do(req)
	utils.HandleError(err)
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	utils.HandleError(err)

	fmt.Fprint(w, string(b))
}

func Start(port int) {
	addr := fmt.Sprintf(":%d", port)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.Use(corsMiddleware)
	router.Use(optionMethodBanMiddleware)

	router.HandleFunc("/", hello)
	router.HandleFunc("/abilities", getabs)

	router.HandleFunc("/signin", auth.SignIn)
	router.HandleFunc("/signup", auth.SignUp)
	router.HandleFunc("/logout", auth.Logout)
	router.HandleFunc("/welcome", auth.Welcome)
	router.HandleFunc("/profile", auth.UserInfo)
	//H4uN
	router.HandleFunc("/makevm", makevm.UploadsHandler)
	router.HandleFunc("/makeqcow", makevm.Makevmfile)
	router.HandleFunc("/isotoqcow2", makevm.Isotoqcow2)
	router.HandleFunc("/listwinvm", makevm.Listwinvm)
	router.HandleFunc("/delqcow2", makevm.Delqcow2)
	router.HandleFunc("/editqcow2", makevm.Editqcow2)
	router.HandleFunc("/startqcow2", makevm.Startqcow2)
	router.HandleFunc("/cloneqcow2", makevm.Cloneqcow2)
	router.HandleFunc("/delwinvm", makevm.Delwinvm)
	router.HandleFunc("/startwinvm", makevm.Startwinvm)
	router.HandleFunc("/resumewinvm", makevm.Resumewinvm)
	router.HandleFunc("/suspendwinvm", makevm.Suspendwinvm)
	router.HandleFunc("/accessvncwindows", makevm.AccessVNCWindows)
	router.HandleFunc("/qcowlist", makevm.Qcowlist)
	router.HandleFunc("/editiso", makevm.EditISO)
	router.HandleFunc("/deliso", makevm.DelISO)
	router.HandleFunc("/makedocker", makevm.Makedocker)
	router.HandleFunc("/dockerlist", makevm.Dockerlist)
	router.HandleFunc("/dockerimagelist", makevm.Dockerimagelist)
	router.HandleFunc("/dockerdestroy", makevm.DockerDestroy)
	router.HandleFunc("/makedockerimage", makevm.MakeDockerImage)
	router.HandleFunc("/destroydockerimage", makevm.DestoryDockerImage)
	router.HandleFunc("/editdockerimage", makevm.EditDockerImage)
	router.HandleFunc("/accessterminaldocker", makevm.AccessTerminalDocker)
	router.HandleFunc("/getwindowslist", makevm.GetWindowslist)
	router.HandleFunc("/getlinuxlist", makevm.GetLinuxlist)
	router.HandleFunc("/accesswindows", makevm.AccessWindows)
	router.HandleFunc("/accesslinux", makevm.AccessLinux)
	router.HandleFunc("/operation_start_linux", makevm.Operation_Start_Linux)
	router.HandleFunc("/operation_start_windows", makevm.Operation_Start_Windows)

	router.HandleFunc("/training", training.Training)
	router.HandleFunc("/trainingcheck", training.ChallengeCheck)
	router.HandleFunc("/createtraining", training.CreateTraining)
	router.HandleFunc("/gettraining", training.GetTraining)
	router.HandleFunc("/edittraining", training.EditTraining)
	router.HandleFunc("/deletetraining", training.DeleteTraining)
	router.HandleFunc("/gettrainings", training.GetAllTrainings)

	router.HandleFunc("/dashboard", dashboard.Dashboard)
	router.HandleFunc("/dashboardbyuser", dashboard.DashboardByUser)

	router.HandleFunc("/dashboardir", caldera.SocketEndpoint)

	router.HandleFunc("/docker", create.DockerRun)
	router.HandleFunc("/vagrant", create.VagrantRun)

	router.HandleFunc("/notification", board.Notification)
	router.HandleFunc("/noticreate", board.NotiCreate)
	router.HandleFunc("/notiedit", board.NotiEdit)

	router.HandleFunc("/scorelist", scoreboard.GetScore)
	router.HandleFunc("/scoremodal", scoreboard.GetScoreModal)
	router.HandleFunc("/scoregraph", scoreboard.GetGraphData)
	log.Fatal(http.ListenAndServe(addr, router))
}

func main() {
	var port int
	scoreboard.Cal()
	fmt.Printf("사용할 포트 입력 (수정 : 3000, 성현 : 8000) : ")
	fmt.Scanf("%d", &port)
	Start(port)
}
