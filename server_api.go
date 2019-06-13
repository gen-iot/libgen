package libgen

import (
	"gitee.com/Puietel/std"
	"time"
)

type Server struct {
}

func NewServer() *Server {
	s := new(Server)
	return s
}

func (this *Server) OnDataRead(buf std.ReadableBuffer, api API) {
	for {
		msg, err := Decode(buf, 1024*1024*1)
		if err != nil {
			if err == ErrNeedMore {
				break
			}
			api.OnError(err)
			return
		}
		this.dispatchInvoke(msg, api)
	}
}

func (this *Server) dispatchInvoke(msg *IOMsg, api API) {
	var out interface{} = nil
	switch GenCommand(msg.Cmd) {
	case CmdDeclareDeviceModel:
		out = api.FetchDevices()
	}

	if out == nil {
		return
	}
	bytes, err := Encode(msg.Id, msg.Cmd, msg.Format, out)
	if err != nil {
		api.OnError(err)
		return
	}
	err = api.Send(time.Second*5, bytes)
	if err != nil {
		api.OnError(err)
	}
}
