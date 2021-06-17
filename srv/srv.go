package srv

import (
	"livego/configure"
	"livego/protocol/rtmp"
	"os"
	e "public/entities"
	conf "public/libs_go/servicelib/config"
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
	return globals.ClickAddress
}

func (s *Service) Metrics() (data []byte) {
	// data, err := json.Marshal(loginAccount)
	// if err != nil {
	// 	return make([]byte, 0)
	// }
	data = make([]byte, 0)
	return
}

func (s *Service) Ext() (data interface{}) {
	data = "Without Ext"
	return
}

func (s *Service) OnVersion(data []byte) {
	// data 为版本实体的序列化数据
	moduleVersion.CheckVersion(data)
	return
}

func (s *Service) Start(c *conf.Setting) (port int, err error) {
	logs.Info("版本：", version)
	stream := rtmp.NewRtmpStream()
	startAPI(stream)

	startRtmp(stream, nil)
	return
}

func (s *Service) Stop() {
	// 解除注册
	time.Sleep(300 * time.Millisecond)
	os.Exit(0)
	return
}

func NewConfig() *conf.Config {
	return &conf.Config{
		EtcdEndpoints:  configure.Config.GetStringSlice("etcdendpoints"),
		ServiceType:    serviceType,
		MonitorService: true,
		Version:        version,
		// LogPath:        "logs",                 // default 'logs'
	}
}
