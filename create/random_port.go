package create

import (
	"bytes"
	"math/rand"
	"os/exec"
	"strconv"
	"time"
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
