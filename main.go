package main

import (
	"os"
	"os/signal"
	"public/libs_go/servicelib"
	"syscall"
	"time"

	"livego/srv"

	"github.com/astaxie/beego/logs"
)

var VERSION = "master"

func main() {
	logs.SetLogger("console")
	defer func() {
		if r := recover(); r != nil {
			logs.Error("livego panic: ", r)
			time.Sleep(1 * time.Second)
		}
		go servicelib.Stop(srv.Ser)
		time.Sleep(300 * time.Millisecond)
	}()

	logs.Info(`
     _     _            ____       
    | |   (_)_   _____ / ___| ___  
    | |   | \ \ / / _ \ |  _ / _ \ 
    | |___| |\ V /  __/ |_| | (_) |
    |_____|_| \_/ \___|\____|\___/ 
	`)

	srv.Scfg = srv.NewConfig()
	err := servicelib.WatchPreSetting(srv.Ser, srv.Scfg)
	if err != nil {
		logs.Error(err)
		return
	}
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	_ = <-chSig
}
