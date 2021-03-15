package main

import (
	"public/libs_go/servicelib"
	"os"
	"os/signal"
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
        version: %s
	`, VERSION)

	srv.Scfg = srv.NewConfig()
	servicelib.WatchPreSetting(srv.Ser, srv.Scfg)
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	_ = <-chSig
}
