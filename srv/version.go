package srv

import (
	"encoding/json"
	e "public/entities"
	"public/models/command"
	"public/libs_go/darwinfslib"
	"public/libs_go/gateway/service"
	"public/libs_go/servicelib"
	"public/libs_go/socketlib"
	"os"
	"time"

	"github.com/astaxie/beego/logs"
)

var moduleVersion = new(ModuleVersion)
var darwinfsProxy = new(darwinfslib.Darwinfs)

type ModuleVersion struct {
	armNode string
	amdNode string
}

//Handle Handle
func (m *ModuleVersion) handle(packet *socketlib.Packet) {
	switch packet.Method {
	case command.SendVersion:
		m.SendVersion(packet)
	}
}

//SendVersion SendVersion
func (m *ModuleVersion) SendVersion(packet *socketlib.Packet) {
	version := new(e.Version)
	err := json.Unmarshal(packet.Content, version)
	if err == nil {
		versionPath := "http://" + globals.GateWay + "/v1/file/" + service.DefaultRegion + loginAccount.LastVersion.Path
		logs.Info("发现新版本", loginAccount.LastVersion.Code, versionPath)
		err = darwinfsProxy.DownloadFile(versionPath, "update.zip")
		if err == nil {
			logs.Info("准备升级,程序退出")
			close()
		} else {
			logs.Error("升级失败", err)
		}
	} else {
		logs.Error("升级指令异常", err)
	}
}

//CheckVersion CheckVersion
func (m *ModuleVersion) CheckVersion() {
	if loginAccount.LastVersion != nil && globals.GateWay != "" {
		versionPath := "http://" + globals.GateWay + "/v1/file/" + service.DefaultRegion + loginAccount.LastVersion.Path
		logs.Info("发现新版本", loginAccount.LastVersion.Code, versionPath)
		err := darwinfsProxy.DownloadFile(versionPath, "update.zip")
		if err == nil {
			logs.Info("准备升级,程序退出")
			close()
		} else {
			logs.Error("升级失败", err)
		}
	}
}

func close() {
	go servicelib.Stop(Ser)
	logs.Info("关闭程序")
	os.Exit(0)
	time.Sleep(300 * time.Millisecond)
}
