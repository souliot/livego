package srv

import (
	"livego/configure"
	"livego/protocol/rtmp"
	"os"
	e "public/entities"
	"public/libs_go/gateway/master"
	conf "public/libs_go/servicelib/config"
	"time"

	"public/libs_go/logs"
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

func (s *Service) SerSetting(data []byte) *master.SerSetting {
	return nil
}

func (s *Service) AppSetting(data []byte) *master.AppSetting {
	return nil
}

func (s *Service) GlobalSetting() interface{} {
	return &e.ServerSetting{}
}

func (s *Service) OnGlobalSetting(c interface{}) {
	globals = c.(*e.ServerSetting)
	etcdEndPoints := configure.Config.GetStringSlice("etcdendpoints")
	initEtcd(etcdEndPoints)
	return
}

func (s *Service) Metrics() (data []byte) {
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

func (s *Service) OnController(data *master.ControllerValue) {
	// data 为控制命令
	close()
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
		EtcdEndpoints: configure.Config.GetStringSlice("etcdendpoints"),
		ServiceType:   serviceType,
		MetricsType:   master.MetricsTypeNone,
		Version:       version,
	}
}
