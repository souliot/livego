package rtmprelay

import (
	"bytes"
	"fmt"
	"io"

	"livego/protocol/amf"
	"livego/protocol/rtmp/core"

	"public/libs_go/logs"
)

var (
	STOP_CTRL = "RTMPRELAY_STOP"
)

type RtmpRelay struct {
	PlayUrl              string
	PublishUrl           string
	cs_chan              chan core.ChunkStream
	sndctrl_chan         chan string
	connectPlayClient    *core.ConnClient
	connectPublishClient *core.ConnClient
	startflag            bool
}

func NewRtmpRelay(playurl *string, publishurl *string) *RtmpRelay {
	return &RtmpRelay{
		PlayUrl:              *playurl,
		PublishUrl:           *publishurl,
		cs_chan:              make(chan core.ChunkStream, 500),
		sndctrl_chan:         make(chan string),
		connectPlayClient:    nil,
		connectPublishClient: nil,
		startflag:            false,
	}
}

func (self *RtmpRelay) rcvPlayChunkStream() {
	logs.Debug("rcvPlayRtmpMediaPacket connectClient.Read...")
	for {
		var rc core.ChunkStream

		if self.startflag == false {
			self.connectPlayClient.Close(nil)
			logs.Debug("rcvPlayChunkStream close: playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
			break
		}
		err := self.connectPlayClient.Read(&rc)

		if err != nil && err == io.EOF {
			break
		}
		//logs.Debug("connectPlayClient.Read return rc.TypeID=%v length=%d, err=%v", rc.TypeID, len(rc.Data), err)
		switch rc.TypeID {
		case 20, 17:
			r := bytes.NewReader(rc.Data)
			vs, err := self.connectPlayClient.DecodeBatch(r, amf.AMF0)

			logs.Debug("rcvPlayRtmpMediaPacket: vs=%v, err=%v", vs, err)
		case 18:
			logs.Debug("rcvPlayRtmpMediaPacket: metadata....")
		case 8, 9:
			self.cs_chan <- rc
		}
	}
}

func (self *RtmpRelay) sendPublishChunkStream() {
	for {
		select {
		case rc := <-self.cs_chan:
			//logs.Debug("sendPublishChunkStream: rc.TypeID=%v length=%d", rc.TypeID, len(rc.Data))
			self.connectPublishClient.Write(rc)
		case ctrlcmd := <-self.sndctrl_chan:
			if ctrlcmd == STOP_CTRL {
				self.connectPublishClient.Close(nil)
				logs.Debug("sendPublishChunkStream close: playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
				break
			}
		}
	}
}

func (self *RtmpRelay) Start() error {
	if self.startflag {
		return fmt.Errorf("The rtmprelay already started, playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
	}

	self.connectPlayClient = core.NewConnClient()
	self.connectPublishClient = core.NewConnClient()

	logs.Debug("play server addr:%v starting....", self.PlayUrl)
	err := self.connectPlayClient.Start(self.PlayUrl, "play")
	if err != nil {
		logs.Debug("connectPlayClient.Start url=%v error", self.PlayUrl)
		return err
	}

	logs.Debug("publish server addr:%v starting....", self.PublishUrl)
	err = self.connectPublishClient.Start(self.PublishUrl, "publish")
	if err != nil {
		logs.Debug("connectPublishClient.Start url=%v error", self.PublishUrl)
		self.connectPlayClient.Close(nil)
		return err
	}

	self.startflag = true
	go self.rcvPlayChunkStream()
	go self.sendPublishChunkStream()

	return nil
}

func (self *RtmpRelay) Stop() {
	if !self.startflag {
		logs.Debug("The rtmprelay already stoped, playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
		return
	}

	self.startflag = false
	self.sndctrl_chan <- STOP_CTRL
}
