package srv

import (
	e "public/entities"
)

var (
	apps    = new(AppSetting)
	globals = new(e.ServerSetting)
)

type AppSetting struct {
	AppName  string `json:"AppName"`
	HTTPPort int    `json:"HttpPort"`
}
