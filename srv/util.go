package srv

import (
	"livego/configure"
	e "public/entities"
	"public/libs_go/gateway/master"
	"runtime"
	"strconv"
	"time"

	"public/libs_go/servicelib"

	"public/libs_go/logs"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	version             = "5.1.1.0"
	serviceType         = strconv.Itoa(e.AccountTypeRtmp)
	accountId           string
	isGetBusinessServer = false
	client              *clientv3.Client
)

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
	accountId = id
	if configure.Config.GetString("version") != "" {
		version = configure.Config.GetString("version")
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
