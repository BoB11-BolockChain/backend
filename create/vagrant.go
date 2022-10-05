package create

import (
	"net/http"

	"github.com/bmatcuk/go-vagrant"
)

func vagrantrun() {
	client, err := vagrant.NewVagrantClient("/home/ar/vagrant_project/test_vagrant")
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
}

func VagrantRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	vagrantrun()
}
