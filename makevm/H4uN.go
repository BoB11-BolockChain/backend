package makevm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/huandu/xstrings"
	"github.com/libvirt/libvirt-go"
)

func findExt(fileName string) string {
	var copyExtensions string
	flag := 1
	for index := len(fileName) - 1; index >= 0; index-- {
		if fileName[index] == '.' {
			flag = 0
			break
		}
		copyExtensions += string(fileName[index])
	}
	if flag == 1 { //확장자 없는 경우 에러 처리
		return ""
	}
	copyExtensions = xstrings.Reverse(copyExtensions)

	return copyExtensions
}

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	uploadFile, header, err := r.FormFile("upload_file") // id가 upload_file이다.
	if err != nil {                                      // 에러 제어
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	// defer는 함수가 종료되기 직전에 실행됨
	defer uploadFile.Close() // 파일을 만들고 닫아줘야함(os자원이라 반납해야함)

	if findExt(header.Filename) == "iso" { //windows .iso파일 업로드
		dirname := "./makevm/uploads/Windows"
		os.MkdirAll(dirname, 0777)                                 // dirname 폴더가 없으면 만들어줌, 777 -> read,write,execute 가능
		filepath := fmt.Sprintf("%s/%s", dirname, header.Filename) // 폴더명/파일명, 파일명은 header에 들어있다.
		file, err := os.Create(filepath)                           // 비어있는 새로운 파일을 만듬

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		io.Copy(file, uploadFile)    // 비어있는 파일에 uploadFile을 복사해준다.
		w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
		fmt.Fprint(w, filepath)      // 어디에 업로드되는지 출력
		defer file.Close()           // 파일을 만들고 닫아줘야함(os자원이라 반납해야함)

	} else if findExt(header.Filename) == "tar" { //Linux 도커 이미지 .tar 업로드
		dirname := "./makevm/uploads/Linux"
		os.MkdirAll(dirname, 0777)                                 // dirname 폴더가 없으면 만들어줌, 777 -> read,write,execute 가능
		filepath := fmt.Sprintf("%s/%s", dirname, header.Filename) // 폴더명/파일명, 파일명은 header에 들어있다.
		file, err := os.Create(filepath)                           // 비어있는 새로운 파일을 만듬

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		io.Copy(file, uploadFile)    // 비어있는 파일에 uploadFile을 복사해준다.
		w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
		fmt.Fprint(w, filepath)      // 어디에 업로드되는지 출력
		defer file.Close()           // 파일을 만들고 닫아줘야함(os자원이라 반납해야함)
	} else {
		fmt.Fprint(w, "확장자를 확인해주세요")
		return
	}
}

// iso list
type FileNamelist struct {
	Num      int    `json:"num"`
	FileName string `json:"title"`
}

func Makevmfile(w http.ResponseWriter, r *http.Request) {
	targetDir := "/var/www/html/back/backend/makevm/uploads/Windows/"
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}

	// var list_ret FileNamelist
	var list_retn []FileNamelist

	for i, file := range files {
		// 파일명
		//fmt.Println(file.Name())
		var list_ret FileNamelist
		list_ret.Num = i + 1
		list_ret.FileName = file.Name()
		list_retn = append(list_retn, list_ret)
		//   // 파일의 절대경로
		//   fmt.Println(fmt.Sprintf("%v/%v", targetDir, file.Name()))
	}
	data := struct {
		Data []FileNamelist `json:"data"`
	}{list_retn}
	json.NewEncoder(w).Encode(data)
	fmt.Println(data)
	w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	fmt.Println("complete vm list")
}

// isExistFile
func isExistFile(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	return true
}

func fileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func ExcuteCMD(script string, arg ...string) {
	cmd := exec.Command(script, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println((err))
	} else {
		fmt.Println(string(output))
	}
}

// makevm list
type Isoname struct {
	ISO_name string
}

// iso to vm(qcow2)
func Isotoqcow2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var isoname Isoname
	// enc := json.NewEncoder(w)
	json.NewDecoder(r.Body).Decode(&isoname)

	TemplateDIR := "/home/ar/user_windows/template/"
	// OriginalDIR := "/home/ar/user_windows/user"
	UploadDIR := "/var/www/html/back/backend/makevm/uploads/Windows/"
	NewVM := TemplateDIR + fileNameWithoutExtension(isoname.ISO_name) + ".qcow2"

	if !isExistFile(TemplateDIR + "virtio-win-0.1.171.iso") {
		ExcuteCMD("wget", "https://fedorapeople.org/groups/virt/virtio-win/direct-downloads/archive-virtio/virtio-win-0.1.171-1/virtio-win-0.1.171.iso", "-P", TemplateDIR)
	}

	// fmt.Println(NewVM)
	// Comp_NewVM:=OriginalDIR+fileNameWithoutExtension(isoname.ISO_name)+"_comp.qcow2"

	//빈 qcow2 파일 생성
	// fmt.Println("1번")
	ExcuteCMD("qemu-img", "create", "-f", "qcow2", NewVM, "30G")

	// iso -> qcow2 생성
	// fmt.Println("2번")
	ExcuteCMD("virt-install", "--name="+fileNameWithoutExtension(isoname.ISO_name), "--ram=4096", "--cpu=host", "--vcpus=1", "--os-type=windows", "--os-variant=win10", "--disk", "path="+NewVM+",format=qcow2,bus=virtio", "--disk", UploadDIR+isoname.ISO_name+",device=cdrom", "--disk", TemplateDIR+"virtio-win-0.1.171.iso,device=cdrom", "--network", "network=default,model=virtio", "--graphics", "vnc,password=pdxf,listen=0.0.0.0", "--import", "--wait", "0", "--check", "all=off", "--cdrom="+UploadDIR+isoname.ISO_name)

}

// makevm eidit list
type EditISOname struct {
	Windows_ISO_NewName string
	ISO_name            string
}

func EditISO(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var isoeditname EditISOname
	json.NewDecoder(r.Body).Decode(&isoeditname)
	UploadDIR := "/var/www/html/back/backend/makevm/uploads/Windows/"
	NEW_ISO_Routes := UploadDIR + isoeditname.Windows_ISO_NewName
	ISO_route := UploadDIR + isoeditname.ISO_name
	ExcuteCMD("mv", "-f", ISO_route, NEW_ISO_Routes)
}

func DelISO(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var isoname Isoname
	json.NewDecoder(r.Body).Decode(&isoname)
	UploadDIR := "/var/www/html/back/backend/makevm/uploads/Windows/"
	ISO_route := UploadDIR + isoname.ISO_name
	ExcuteCMD("rm", "-rf", ISO_route)
}

// virsh list
type Vrishlist struct {
	ID     uint
	Domain string
	State  string
	Port   int
}

func GetStatemsg(msg int) string {
	retdata := ""
	switch msg {
	case 0:
		retdata = "NOSTATE"
	case 1:
		retdata = "RUNNING"
	case 2:
		retdata = "BLOCKED"
	case 3:
		retdata = "PAUSED"
	case 4:
		retdata = "SHUTDOWN"
	case 5:
		retdata = "SHUTOFF"
	case 6:
		retdata = "CRASHED"
	case 7:
		retdata = "PMSUSPENDED"
	default:
		retdata = "NOSTATE"
	}
	return retdata
}

func Listwinvm(w http.ResponseWriter, r *http.Request) { //Windows list --all

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		println(err)
		return
	}
	defer conn.Close()

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		println(err)
		return
	}

	var retrunvirshlist []Vrishlist
	for _, dom := range doms {
		var tempvirshlist Vrishlist
		name, err := dom.GetName()
		if err != nil {
			println(err)
			return
		}
		dominfo, err := dom.GetInfo()
		if err != nil {
			println(err)
			return
		}
		domid, err := dom.GetID()
		if err != nil {
			println(err)
			return
		}

		tempvirshlist.Domain = name
		tempvirshlist.ID = domid
		tempvirshlist.State = GetStatemsg(int(dominfo.State))

		cmd := exec.Command("virsh", "vncdisplay", name)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			fmt.Println((err))
		}
		tempvncport := string(output)
		r, _ := regexp.Compile("[0-9]")
		vncport := r.FindString(tempvncport)
		returnvncport, err := strconv.Atoi(vncport)
		if err != nil {
			fmt.Println((err))
		}
		tempvirshlist.Port = returnvncport
		// fmt.Println(name)
		// fmt.Println(domid)
		// fmt.Println(GetStatemsg(int(dominfo.State)))
		retrunvirshlist = append(retrunvirshlist, tempvirshlist)
		dom.Free()
	}

	doms_in, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		println(err)
		return
	}

	for _, dom := range doms_in {
		var tempvirshlist Vrishlist
		name, err := dom.GetName()
		if err != nil {
			println(err)
			return
		}
		dominfo, err := dom.GetInfo()
		if err != nil {
			println(err)
			return
		}
		tempvirshlist.Domain = name
		tempvirshlist.State = GetStatemsg(int(dominfo.State))
		tempvirshlist.Port = 4444
		// fmt.Println(name)
		// fmt.Println(GetStatemsg(int(dominfo.State)))
		retrunvirshlist = append(retrunvirshlist, tempvirshlist)
		fmt.Println(retrunvirshlist)
		dom.Free()
	}

	sort.Slice(retrunvirshlist, func(i, j int) bool {
		return retrunvirshlist[i].Port < retrunvirshlist[j].Port
	})

	for i := 0; i < len(retrunvirshlist); i++ {
		if retrunvirshlist[i].Port == 4444 {
			retrunvirshlist[i].Port = 0
			continue
		}
		retrunvirshlist[i].Port = retrunvirshlist[i].Port + 5900
	}

	data := struct {
		Data []Vrishlist `json:"data"`
	}{retrunvirshlist}
	json.NewEncoder(w).Encode(data)
	fmt.Println(data)
	w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
}

func Qcowlist(w http.ResponseWriter, r *http.Request) {
	targetDir := "/home/ar/user_windows/template"
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		panic(err)
	}

	// var list_ret FileNamelist
	var list_retn []FileNamelist

	for i, file := range files {
		// 파일명
		//fmt.Println(file.Name())
		if findExt(file.Name()) == "qcow2" {
			var list_ret FileNamelist
			list_ret.Num = i + 1
			list_ret.FileName = file.Name()
			list_retn = append(list_retn, list_ret)
			//   // 파일의 절대경로
			//   fmt.Println(fmt.Sprintf("%v/%v", targetDir, file.Name()))
		}
	}
	data := struct {
		Data []FileNamelist `json:"data"`
	}{list_retn}
	json.NewEncoder(w).Encode(data)
	fmt.Println(data)
	w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	fmt.Println("complete qcow2 list")
}

type Qcow2file struct {
	Newfilename    string
	Originfilename string
}

func Editqcow2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var qcowname Qcow2file
	json.NewDecoder(r.Body).Decode(&qcowname)
	TemplateDIR := "/home/ar/user_windows/template/"
	NEW_qcow2name := TemplateDIR + qcowname.Newfilename
	Origin_qcow2name := TemplateDIR + qcowname.Originfilename
	ExcuteCMD("mv", "-f", Origin_qcow2name, NEW_qcow2name)
}

type Qcow2file2 struct {
	Filename string
}

func Delqcow2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var qcowname Qcow2file2
	json.NewDecoder(r.Body).Decode(&qcowname)
	TemplateDIR := "/home/ar/user_windows/template/"
	Qcow2_route := TemplateDIR + qcowname.Filename
	ExcuteCMD("rm", "-rf", Qcow2_route)
}

//나중에 vnc 연결로 버튼 빼자
// //포트 반환
// cmd := exec.Command("virsh", "vncdisplay", fileNameWithoutExtension(isoname.ISO_name))
// output, err := cmd.CombinedOutput()
// if err != nil {
// 	fmt.Println((err))
// }
// // } else {
// // 	fmt.Println(string(output))
// // }
// temp_port_num := string(output)
// // fmt.Println("asdf" + temp_port_num[1:2] + "asdf")
// temp_num, err := strconv.Atoi(temp_port_num[1:2])
// fmt.Println(temp_num, err)
// port_num := 5900 + temp_num
// // fmt.Println(port_num)
// isodata := struct {
// 	ISO_port int    `json:"ISO__Port"`
// 	ISO_name string `json:"ISO__Name"`
// }{port_num, isoname.ISO_name}
// enc.Encode(isodata)
// // vm 끄기
// fmt.Println("3번")
// ExcuteCMD("virsh", "shutdown", fileNameWithoutExtension(isoname.ISO_name))
// Sleep()
// // qcow2 이미지 압축
// fmt.Println("4번")
// ExcuteCMD("qemu-img","convert" ,"-O" ,"qcow2" ,"-c" ,NewVM, Comp_NewVM)

// ExcuteCMD("virt-install", "--name="+fileNameWithoutExtension(isoname.ISO_name), "--ram=4096", "--cpu=host", "--vcpus=1" ,"--os-type=windows" ,"--os-variant=win10" ,"--disk path="+NewVM+",device=disk,bus=virtio,format=qcow2","--network network=default,model=virtio", "--graphics vnc,password=test,listen=0.0.0.0", "--import", "--wait 0" ,"--check all=off")
// fmt.Println("5번")
// }

// 이거 나중에 vm 실행할때 쓰자
// func Isotoqcow2(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 	r.ParseForm()
// 	var isoname Isoname
// 	json.NewDecoder(r.Body).Decode(&isoname)
// 	// fmt.Println(isoname)
// 	// OriginalDIR:="/home/ar/iso"
// 	CreateDIR:="/home/ar/iso/VM/" + fileNameWithoutExtension(isoname.ISO_name)
// 	// fmt.Println(CreateDIR)
// 	if !isExistFile(CreateDIR){
// 		// ExcuteCMD("mkdir", CreateDIR)
// 		// NewVM:=CreateDIR+fileNameWithoutExtension(isoname.ISO_name)+".qcow2"
// 		ExcuteCMD("virt-install", "--name="+fileNameWithoutExtension(isoname.ISO_name), "--ram=4096", "--cpu=host", "--vcpus=1" ,"--os-type=windows" ,"--os-variant=win10" ,"--disk path="+NewVM+",device=disk,bus=virtio,format=qcow2","--network network=default,model=virtio", "--graphics vnc,password=test,listen=0.0.0.0", "--import", "--wait 0" ,"--check all=off")
// 		// fmt.Println(CreateDIR)
// 	}
// }

type Qcow2startstruct struct {
	Filename string
}

func Startqcow2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var filename Qcow2startstruct
	json.NewDecoder(r.Body).Decode(&filename)
	TemplateDIR := "/home/ar/user_windows/template/"
	NewDIR := TemplateDIR + filename.Filename
	ExcuteCMD("virt-install", "--name="+fileNameWithoutExtension(filename.Filename), "--ram=4096", "--cpu=host", "--vcpus=1", "--os-type=windows", "--os-variant=win10", "--disk", "path="+NewDIR+",device=disk,bus=virtio,format=qcow2", "--network", "network=default,model=virtio", "--graphics", "vnc,password=pdxf,listen=0.0.0.0", "--import", "--wait", "0", "--check", "all=off")
}

type VNCWindows struct {
	VMname   string
	VNC_port string
}

func AccessVNCWindows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var vncwindows VNCWindows
	json.NewDecoder(r.Body).Decode(&vncwindows)
	tempnum1, err := strconv.Atoi(vncwindows.VNC_port)
	if err == nil {
		fmt.Println("Atoi Err")
	}
	tempnum2 := tempnum1 + 180
	restrnum := strconv.Itoa(tempnum2)
	ExcuteCMD("sudo", "sh", "/usr/share/novnc/utils/launch.sh", "--vnc", "localhost:"+vncwindows.VNC_port, "--ssl-only", "--listen", restrnum)
}

type Qcow2clonestruct struct {
	Filename      string
	Clonefilename string
}

func Cloneqcow2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var filename Qcow2clonestruct
	json.NewDecoder(r.Body).Decode(&filename)

	TemplateDIR := "/home/ar/user_windows/template/"
	OriginDIR := TemplateDIR + filename.Filename
	CloneDIR := TemplateDIR + filename.Clonefilename
	ExcuteCMD("cp", "-rf", OriginDIR, CloneDIR)
}

// makedocker
type Docker_ struct {
	Docker_name string
}

func Makedocker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var do_name Docker_
	json.NewDecoder(r.Body).Decode(&do_name)
	// fmt.Println(do_name.Docker_name)
	ExcuteCMD("docker", "pull", do_name.Docker_name)
	fmt.Println("도커 성공적")
}

func AccessTerminalDocker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var tempdata DockerTemp
	json.NewDecoder(r.Body).Decode(&tempdata)
	ExcuteCMD("gotty", "--once", "-w", "docker", "run", "-it", "--rm", tempdata.Image_ID, "/bin/bash")
}

// 커스텀성현

type DockerResponse struct {
	Ssh  string `json:"ssh"`
	Http string `json:"http"`
}

func dockerun(port1, port2, linux_name string) int {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	// fmt.Println("thisis" + linux_name)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: linux_name,
		//apt update; apt install net-tools openssh-server vim -y; mkdir /var/run/sshd; echo 'root:root' |chpasswd; sed -ri 's/^#?PermitRootLogin\s+.*/PermitRootLogin yes/' /etc/ssh/sshd_config; sed -ri 's/UsePAM yes/#UsePAM yes/g' /etc/ssh/sshd_config; mkdir /root/.ssh; service ssh start; sleep infinity
		//Cmd: []string{"apt", "update;", "apt", "install", "net-tools", "openssh-server", "vim", "-y;", "mkdir", "/var/run/sshd;", "echo", "'root:root'", "|chpasswd;", "sed", "-ri", "'s/^#?PermitRootLogin\\s+.*/PermitRootLogin", "yes/'", "/etc/ssh/sshd_config;", "sed", "-ri", "'s/UsePAM", "yes/#UsePAM", "yes/g'", "/etc/ssh/sshd_config;", "mkdir", "/root/.ssh;", "service", "ssh", "start;", "sleep", "infinity"},
		// Cmd: []string{"sleep", "infinity"},
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

func DockerRun(linux_name string) DockerResponse {
	var portstr1, portstr2 string = RandPort()
	fmt.Println(portstr1)
	fmt.Println(portstr2)
	if dockerun(portstr1, portstr2, linux_name) == 1 {
		resp := DockerResponse{portstr1, portstr2}
		return resp
	}
	return DockerResponse{"0", "0"}
}

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

// Docker list
type Docker_list struct {
	ContainerID string       `json:"ContainerID"`
	IMAGE       string       `json:"IMAGE"`
	STATUS      string       `json:"STATUS"`
	PORTS       []types.Port `json:"PORTS"`
}
type data struct {
	Data []Docker_list `json:"data"`
}

func Dockerlist(w http.ResponseWriter, r *http.Request) {

	// var list_ret FileNamelist
	// var dockerlist_rtn []Docker_list
	dockerlist_rtn := Dockerps()
	json.NewEncoder(w).Encode(dockerlist_rtn)
	// fmt.Println(dockerlist_rtn)
	w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	fmt.Println("docker list complete")
}

func Dockerps() data {

	var dockerlist_retn []Docker_list
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		var dockerlist_rtn Docker_list
		// fmt.Println(container.ID)
		// fmt.Println(container.Image)
		// fmt.Println(container.Status)
		// fmt.Println(container.Ports)
		dockerlist_rtn.ContainerID = container.ID
		dockerlist_rtn.ContainerID = dockerlist_rtn.ContainerID[:15]
		dockerlist_rtn.IMAGE = container.Image
		dockerlist_rtn.STATUS = container.Status
		dockerlist_rtn.PORTS = container.Ports
		dockerlist_retn = append(dockerlist_retn, dockerlist_rtn)
	}
	data := struct {
		Data []Docker_list `json:"data"`
	}{dockerlist_retn}

	// fmt.Println(data)
	return data
}

// dockerlist
type Dockerlistresp struct {
	Docker_image       string
	Docker_containerID string
}

func DockerDestroy(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var _name Dockerlistresp
	json.NewDecoder(r.Body).Decode(&_name)
	// fmt.Println(_name.Docker_containerID)
	ExcuteCMD("docker", "rm", "-f", _name.Docker_containerID)
}

func Dockerimagelist(w http.ResponseWriter, r *http.Request) {
	dockerimagelist_rtn := Dockerimageps()
	// fmt.Println(dockerimagelist_rtn)
	json.NewEncoder(w).Encode(dockerimagelist_rtn)
	fmt.Println(dockerimagelist_rtn)
	w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	fmt.Println("docker image list complete")
}

// Docker image list
type Docker_image_list struct {
	REPOSITORY string `json:"REPOSITORY"`
	IMAGE_ID   string `json:"IMAGE_ID"`
	CREATED    int    `json:"CREATED"`
}

type data1 struct {
	Data []Docker_image_list `json:"data"`
}

func Dockerimageps() data1 {
	var dockerimagelist_retn []Docker_image_list
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		var dockerimagelist_rtn Docker_image_list
		// fmt.Println(image.RepoTags)
		// fmt.Println(image.ID)
		// fmt.Println(image.Created)
		// fmt.Println(reflect.TypeOf(image.RepoTags))
		if image.RepoTags != nil {
			for _, temp := range image.RepoTags {
				dockerimagelist_rtn.REPOSITORY = temp
				dockerimagelist_rtn.IMAGE_ID = image.ID[7:21]
				dockerimagelist_rtn.CREATED = int(image.Created)
				dockerimagelist_retn = append(dockerimagelist_retn, dockerimagelist_rtn)
			}
		}
	}

	data := struct {
		Data []Docker_image_list `json:"data"`
	}{dockerimagelist_retn}

	// fmt.Println(data)
	return data
}

// Win vm del
type WinVM struct {
	WinVMDomain string
}

func Delwinvm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var winvmdata WinVM
	json.NewDecoder(r.Body).Decode(&winvmdata)
	// fmt.Println(winvmdata)
	ExcuteCMD("virsh", "destroy", winvmdata.WinVMDomain)
	ExcuteCMD("virsh", "undefine", winvmdata.WinVMDomain)
}

func Suspendwinvm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var winvmdata WinVM
	json.NewDecoder(r.Body).Decode(&winvmdata)
	// fmt.Println(winvmdata)
	ExcuteCMD("virsh", "suspend", winvmdata.WinVMDomain)
}

func Startwinvm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var winvmdata WinVM
	json.NewDecoder(r.Body).Decode(&winvmdata)
	// fmt.Println(winvmdata)
	ExcuteCMD("virsh", "start", winvmdata.WinVMDomain)
}

func Resumewinvm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var winvmdata WinVM
	json.NewDecoder(r.Body).Decode(&winvmdata)
	// fmt.Println(winvmdata)
	ExcuteCMD("virsh", "resume", winvmdata.WinVMDomain)
}

func MakeDockerImage(w http.ResponseWriter, r *http.Request) { //도커 컨테이너 -> 이미지
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var _name Dockerlistresp
	json.NewDecoder(r.Body).Decode(&_name)
	// fmt.Println(_name)
	ExcuteCMD("docker", "commit", _name.Docker_containerID, _name.Docker_image)
	fmt.Println("도커 이미지 생성 완료")
}

type DockerTemp struct {
	Image_ID       string
	Container_Name string
}

func DestoryDockerImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	var temp DockerTemp
	json.NewDecoder(r.Body).Decode(&temp)
	// makestring := strings.Index(temp.Container_Name, ":")
	// temp.Container_Name = temp.Container_Name[:makestring]

	// code_temp := "$(docker ps -aq --filter ancestor=" + temp.Container_Name + ")"
	// fmt.Println(code_temp)
	// ExcuteCMD("docker", "rm", "-f", code_temp)
	cmd := exec.Command("docker", "ps", "-aq", "--filter", "ancestor="+temp.Image_ID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println((err))
	} else {
		fmt.Println(string(output))
	}

	if len(string(output)) == 0 {
		fmt.Println("이게 없네;;")
		ExcuteCMD("docker", "rmi", "-f", temp.Image_ID)
	} else if string(output)[len(string(output))-1] == '\n' {

		tmp_string := strings.Split(string(output), "\n")
		fmt.Println(tmp_string)
		for i := 0; i < len(tmp_string); i++ {
			string_temp := tmp_string[i][:len(tmp_string[i])-1]
			ExcuteCMD("docker", "rm", "-f", string_temp)
			fmt.Print("??")
		}
	} else {
		ExcuteCMD("docker", "rm", "-f", string(output))
		fmt.Print("!!")
	}
	ExcuteCMD("docker", "rmi", "-f", temp.Container_Name)
	fmt.Println("Docker Image Delete Complete")
}

func EditDockerImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()

	type DockerTemp struct {
		Docker_image      string
		Docker_REPOSITORY string
	}

	var temp DockerTemp
	json.NewDecoder(r.Body).Decode(&temp)
	ExcuteCMD("docker", "image", "tag", temp.Docker_REPOSITORY, temp.Docker_image)
	ExcuteCMD("docker", "rmi", temp.Docker_REPOSITORY)
	fmt.Println("Edit Docker Image Name Complete")
}
