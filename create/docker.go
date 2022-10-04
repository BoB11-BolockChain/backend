package create

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

func RandRange(min, max int) int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	return randNum
}

func RandPort() (string, string) {
	for {
		var port1 int = RandRange(1024, 65534)
		var port2 int = port1 + 1

		var s1 bytes.Buffer
		s1.WriteString("netstat -antul | grep ':")
		s1.WriteString(strconv.Itoa(port1))
		s1.WriteString("'")

		var s2 bytes.Buffer
		s2.WriteString("netstat -antul | grep ':")
		s2.WriteString(strconv.Itoa(port2))
		s2.WriteString("'")

		cmd1 := exec.Command("sh", "-c", s1.String())
		_, err1 := cmd1.Output()

		cmd2 := exec.Command("sh", "-c", s2.String())
		_, err2 := cmd2.Output()

		if err1 != nil {
			if err2 != nil {
				return strconv.Itoa(port1), strconv.Itoa(port2)
			}
		}
	}
}

func dockerun(port1, port2 string) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	reader, err := cli.ImagePull(ctx, "docker.io/library/ubuntu", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "ubuntu",
		Cmd:   []string{"sleep", "infinity"},
		Tty:   false,
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

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

func DockerRun(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	// fmt.Fprint(w, r.Form)

	var portstr1, portstr2 string = RandPort()
	fmt.Println(portstr1)
	fmt.Println(portstr2)
	dockerun(portstr1, portstr2)
}
