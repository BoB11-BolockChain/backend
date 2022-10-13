package create

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/bmatcuk/go-vagrant"
)

type VagrantArgs struct {
	Id string `json:"id"`
}

type VagrantResponse struct {
	Id   string `json:"id"`
	Rdp  string `json:"rdp"`
	Http string `json:"http"`
}

func vagrantrun(route string) int {
	client, err := vagrant.NewVagrantClient(route)
	if err != nil {
		panic(err)
	}
	upcmd := client.Up()
	upcmd.Verbose = true
	if err := upcmd.Run(); err != nil {
		panic(err)
	}
	if upcmd.Error != nil {
		panic(err)
	}
	return 1
}

func statusmapping(id, port string) string {
	route := "/home/vagrant/" + id + "/"
	err := exec.Command("mkdir", route).Run()
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	b.WriteString("/home/vagrant/")
	b.WriteString(id)
	b.WriteString("/Vagrantfile")

	file, err := os.Create(b.String())
	if err != nil {
		panic(err)
	}

	ex_vf, err := os.Open("/home/vagrant/ex_Vagrantfile")
	if err != nil {
		panic(err)
	}

	ex_vf_read := bufio.NewReader(ex_vf)

	for {
		line, _, err := ex_vf_read.ReadLine()
		if err == io.EOF {
			break
		}
		temp_str1 := string(line[:])
		temp_str2 := strings.Replace(temp_str1, ";port;", port, 1)
		file.WriteString(temp_str2)
		file.WriteString("\n")
	}
	return route
}

func VagrantRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method != "POST" {
		return
	}
	r.ParseForm()
	// fmt.Fprint(w, r.Form)

	var args VagrantArgs
	json.NewDecoder(r.Body).Decode(&args)
	fmt.Println(args)

	var id string = args.Id
	var portstr1, portstr2 string = RandPort()
	fmt.Println(portstr1)
	fmt.Println(portstr2)

	if vagrantrun(statusmapping(id, portstr1)) == 1 {
		resp := VagrantResponse{id, portstr1, portstr2}
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		enc.Encode(resp)
	}
}
