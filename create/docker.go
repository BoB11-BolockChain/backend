package create

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

type DockerArgs struct {
	Id string `json:"id"`
}

type DockerResponse struct {
	Ssh  string `json:"ssh"`
	Http string `json:"http"`
}

func dockerun(port1, port2 string) int {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	// reader, err := cli.ImagePull(ctx, "docker.io/library/ubuntu", types.ImagePullOptions{})
	// if err != nil {
	// 	panic(err)
	// }

	// defer reader.Close()
	// io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "s_ubuntu",
		// Cmd:   []string{"sleep", "infinity"},
		Tty: false,
		ExposedPorts: nat.PortSet{
			nat.Port("22/tcp"): {},
			nat.Port("80/tcp"): {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"22/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port1,
				},
			},
			"80/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port2,
				},
			},
		},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	/*
		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}
	*/
	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	return 1
}

func DockerRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method != "POST" {
		return
	}
	r.ParseForm()
	// fmt.Fprint(w, r.Form)
	var args DockerArgs
	json.NewDecoder(r.Body).Decode(&args)
	fmt.Println(args)

	var portstr1, portstr2 string = RandPort()
	fmt.Println(portstr1)
	fmt.Println(portstr2)
	if dockerun(portstr1, portstr2) == 1 {
		resp := DockerResponse{portstr1, portstr2}
		enc := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")
		enc.Encode(resp)
	}
}
