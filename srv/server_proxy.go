package srv

import (
	"encoding/json"
	"livego/configure"
	"os"
	e "public/entities"
	"public/libs_go/gateway/master"
	"public/libs_go/gateway/service"
	"public/libs_go/socketlib"
	"public/libs_go/tcpclientlib"
	"public/models/command"
	"runtime"
	"strconv"
	"time"

	"public/libs_go/servicelib"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/coreos/etcd/clientv3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	serverProxy         = new(ServerProxy)
	sysCommand          = socketlib.NewCommand()
	version             = "5.1.1.0"
	serviceType         = strconv.Itoa(e.AccountTypeRtmp)
	loginAccount        = new(e.Account)
	isGetBusinessServer = false
	client              *clientv3.Client
)

type ServerProxy struct {
	cli      *tcpclientlib.SocketClient
	stop     bool
	isLoging bool
}

//Load Load
func (s *ServerProxy) Load(comServer string) {
	s.cli = new(tcpclientlib.SocketClient)
	s.cli.Load(comServer, s, false)
	s.cli.Start()
	s.stop = false
}

func initEtcd(endpoints []string) {
	var err error
	for {
		client, err = master.GetEtcdClient(endpoints, 5)
		if err != nil {
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
}

func getAccountId() (id string) {
	id = configure.Config.GetString("NodeId")
	if id == "" {
		id = primitive.NewObjectID().Hex()
		configure.Config.Set("NodeId", id)
		err := configure.Config.WriteConfigAs(configure.Config.ConfigFileUsed())
		if err != nil {
			logs.Error(err)
		}
	}
	loginAccount.AccountID = id
	loginAccount.LoginName = id
	loginAccount.LoginPwd = id
	loginAccount.AccountType = e.AccountTypeRtmp
	loginAccount.Platform = getPlatform()
	loginAccount.IP = getIP()
	if configure.Config.GetString("version") != "" {
		version = configure.Config.GetString("version")
	}
	loginAccount.Version = version
	return
}

//MessageReceived MessageReceived
func (s *ServerProxy) MessageReceived(packet *socketlib.Packet) {
	switch {
	case packet.Module == command.ModuleVersion:
		moduleVersion.handle(packet)
	case packet.Module == command.ModuleAccount && packet.Method == command.RestartTerminal:
		beego.Info("收到重启命令：", packet.Module, packet.Method, string(packet.Content))
		servicelib.Stop(Ser)
		time.Sleep(300 * time.Millisecond)
		os.Exit(0)
	}
}

//StateChanged StateChanged
func (s *ServerProxy) StateChanged(err error, socketClient *tcpclientlib.SocketClient) {
	if err == nil {
		if !s.isLoging {
			s.Login()
		}
	} else {
		if !s.stop {
			GetComServer()
		}
	}
}

func (s *ServerProxy) Login() {
	s.isLoging = true
	GetBusinessServer()
	getAccountId()
	b, err := s.cli.SendMessage(sysCommand.ModuleSystem, sysCommand.Login, loginAccount, 30)

	if err != nil {
		logs.Error("登录服务失败", err)
		time.Sleep(3 * time.Second)
		logs.Info("重新登录服务……")
		go s.Login()
		return
	}

	err = json.Unmarshal(b, loginAccount)
	if err != nil {
		logs.Error("登录服务失败", err)
		time.Sleep(3 * time.Second)
		logs.Info("重新登录服务……")
		go s.Login()
		return
	}
	logs.Info("登录成功")
	moduleVersion.CheckVersion()
	s.isLoging = false
	return
}

func (s *ServerProxy) Logoff() {
	if s != nil && s.cli != nil {
		s.cli.SendMessage(command.ModuleSystem, command.Logoff, loginAccount, 30)
	}

	return
}

func (s *ServerProxy) Stop() {
	s.stop = true
	s.Logoff()
	if s != nil && s.cli != nil {
		s.cli.Stop()
	}
}

//CheckVersion CheckVersion
func (s *ServerProxy) CheckVersion(version *e.Version) (*e.Version, error) {
	newVersion := new(e.Version)
	buffer, err := s.cli.SendMessage(command.ModuleVersion, command.CheckVersion, version, 60)
	if err == nil && len(buffer) > 0 {
		err = json.Unmarshal(buffer, newVersion)
	}
	return newVersion, err
}

func GetComServer() {
	serviceType := strconv.FormatInt(e.AccountTypeCom, 10)
	for {
		m, err := master.GetService(service.DefaultRegion, serviceType, client, master.Random)
		if err != nil {
			logs.Error("获取通信服务失败,重新获取中", err)
			time.Sleep(time.Second * 3)
			continue
		}
		logs.Info("获取通讯服务成功！")
		serverProxy.Load(m.Address)
		break
	}
	return
}

func GetBusinessServer() {
	isGetBusinessServer = true
	serviceType := strconv.FormatInt(e.AccountTypeBusiness, 10)
	for {
		_, err := master.GetService(service.DefaultRegion, serviceType, client, master.Random)
		if err != nil {
			logs.Error("获取业务服务失败,重新获取中", err)
			time.Sleep(time.Second * 3)
			continue
		}
		isGetBusinessServer = false
		logs.Info("获取业务服务成功")
		break
	}
	return
}

func getPlatform() int {
	switch runtime.GOARCH {
	case "amd64":
		return e.PlatformX64
	case "arm":
		return e.PlatformARM
	case "386":
		return 2
	default:
		return -1
	}
}

func getIP() string {
	ip, err := servicelib.GetIP()
	if err != nil {
		return "127.0.0.1"
	}
	return ip.String()
}
