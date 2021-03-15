package srv

import (
	"livego/configure"
	"livego/protocol/rtmp"
	e "public/entities"
	"public/libs_go/servicelib"
	conf "public/libs_go/servicelib/config"
	"public/libs_go/servicelib/logcollect"
	"public/libs_go/servicelib/monitor"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
)

var (
	Scfg *conf.Config
	Ser  = &Service{}
)

type Service struct{}

func (s *Service) NodeId() (id string) {
	id = getAccountId()
	return
}

func (s *Service) SaveNodeId(v string) {
	return
}
func (s *Service) AppSetting() interface{} {
	return &AppSetting{}
}
func (s *Service) GlobalSetting() interface{} {
	return &e.ServerSetting{}
}
func (s *Service) UpdateGlobalSetting(c interface{}) (clickaddress string) {
	globals = c.(*e.ServerSetting)
	etcdEndPoints := configure.Config.GetStringSlice("etcdendpoints")
	initEtcd(etcdEndPoints)
	go GetComServer()
	return globals.ClickAddress
}

func (s *Service) Start(c *conf.Setting) (port int, err error) {
	logs.Info("版本：", version)
	stream := rtmp.NewRtmpStream()
	startAPI(stream)

	startRtmp(stream, nil)
	return
}

func (s *Service) Stop() {
	go serverProxy.Stop()
	monitor.StopMonitor()
	// 停止日志收集
	logcollect.StopLog()
	// 解除注册
	servicelib.UnRegister()
	time.Sleep(300 * time.Millisecond)
	os.Exit(0)
	return
}

func NewConfig() *conf.Config {
	return &conf.Config{
		EtcdEndpoints: configure.Config.GetStringSlice("etcdendpoints"),
		ServiceType:   serviceType,
		// ClickAddress:   "tcp://192.168.0.8:9000?username=default&password=watrix888",
		MonitorService: true,
		MonitorSystem:  false,
		LogCollect:     true,
		// LogPath:        "logs",                 // default 'logs'
		// LogCollectPath: []string{"logs/*.log"}, // default LogPath/*.log
	}
}
