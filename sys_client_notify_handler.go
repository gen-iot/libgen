package libgen

import (
	"gitee.com/SuzhenProjects/libgen/rpcx"
	"log"
	"time"
)

type DeviceControlHandler func(req *ControlDeviceRequest) (*ControlDeviceResponse, error)
type DeviceStatusHandler func(notify *DeviceStatusNotify)

var gDeviceControlHandler DeviceControlHandler
var gDeviceStatusHandler DeviceStatusHandler


func onPing(callable rpcx.Callable, req *Ping) (*Pong, error) {
	log.Println("receive ping req.time =", req.Time, " delta is ", time.Now().Sub(req.Time))
	return &Pong{Time: time.Now(), Msg: "client pong"}, nil
}

func onDeviceControl(callable rpcx.Callable, req *ControlDeviceRequest) (*ControlDeviceResponse, error) {
	if gDeviceControlHandler != nil {
		return gDeviceControlHandler(req)
	}
	return new(ControlDeviceResponse), nil
}

func onDeviceStatus(callable rpcx.Callable, notify *DeviceStatusNotify) (*BaseResponse, error) {
	if gDeviceStatusHandler != nil {
		go gDeviceStatusHandler(notify)
	}
	return &BaseResponse{}, nil
}

func RegOnDeviceControlHandler(fn DeviceControlHandler) {
	gDeviceControlHandler = fn
}

func RegOnDeviceStatusHandler(fn DeviceStatusHandler) {
	gDeviceStatusHandler = fn
}
