package srv

import (
	"livego/configure"
	"livego/protocol/api"
	"livego/protocol/hls"
	"livego/protocol/httpflv"
	"livego/protocol/rtmp"
	"net"

	"github.com/astaxie/beego/logs"
)

func startHls() *hls.Server {
	hlsAddr := configure.Config.GetString("hls_addr")
	hlsListen, err := net.Listen("tcp", hlsAddr)
	if err != nil {
		logs.Error(err)
	}

	hlsServer := hls.NewServer()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("HLS server panic: ", r)
			}
		}()
		logs.Info("HLS listen On ", hlsAddr)
		hlsServer.Serve(hlsListen)
	}()
	return hlsServer
}

var rtmpAddr string

func startRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
	rtmpAddr = configure.Config.GetString("rtmp_addr")

	rtmpListen, err := net.Listen("tcp", rtmpAddr)
	if err != nil {
		logs.Error(err)
	}

	var rtmpServer *rtmp.Server

	if hlsServer == nil {
		rtmpServer = rtmp.NewRtmpServer(stream, nil)
		logs.Info("HLS server disable....")
	} else {
		rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
		logs.Info("HLS server enable....")
	}

	defer func() {
		if r := recover(); r != nil {
			logs.Error("RTMP server panic: ", r)
		}
	}()
	logs.Info("RTMP Listen On ", rtmpAddr)
	go rtmpServer.Serve(rtmpListen)
}

func startHTTPFlv(stream *rtmp.RtmpStream) {
	httpflvAddr := configure.Config.GetString("httpflv_addr")

	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		logs.Error(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("HTTP-FLV server panic: ", r)
			}
		}()
		logs.Info("HTTP-FLV listen On ", httpflvAddr)
		hdlServer.Serve(flvListen)
	}()
}

func startAPI(stream *rtmp.RtmpStream) {
	apiAddr := configure.Config.GetString("api_addr")

	if apiAddr != "" {
		opListen, err := net.Listen("tcp", apiAddr)
		if err != nil {
			logs.Error(err)
		}
		opServer := api.NewServer(stream, rtmpAddr)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logs.Error("HTTP-API server panic: ", r)
				}
			}()
			logs.Info("HTTP-API listen On ", apiAddr)
			opServer.Serve(opListen)
		}()
	}
}
