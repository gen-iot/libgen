package libgen

import (
	"errors"
	"gitee.com/gen-iot/rpcx"
	"log"
	"time"
)

type DeviceControlHandler func(req *ControlDeviceRequest) (*ControlDeviceResponse, error)
type DeviceStatusHandler func(notify *DeviceStatusNotify)
type TransportDataHandler func(req *TransportDataRequest) (map[string]interface{}, error)

var gDeviceControlHandler DeviceControlHandler
var gDeviceStatusHandler DeviceStatusHandler
var gDataTransportHandler TransportDataHandler

func pong(ctx rpcx.Context, req *Ping) (*Pong, error) {
	log.Println("receive ping req.time =", req.Time, " delta is ", time.Now().Sub(req.Time))
	return &Pong{Time: time.Now(), Msg: "client pong"}, nil
}

var errAppNotImpControl = errors.New("app not support control device yet")

func onDeviceControl(ctx rpcx.Context, req *ControlDeviceRequest) (*ControlDeviceResponse, error) {
	if gDeviceControlHandler != nil {
		return gDeviceControlHandler(req)
	}
	return nil, errAppNotImpControl
}

func onDeviceStatusDelivery(ctx rpcx.Context, notify *DeviceStatusNotify) (*BaseResponse, error) {
	if gDeviceStatusHandler != nil {
		go gDeviceStatusHandler(notify)
	}
	return &BaseResponse{}, nil
}

func onDataTransport(ctx rpcx.Context, req *TransportDataRequest) (map[string]interface{}, error) {
	if gDataTransportHandler != nil {
		return gDataTransportHandler(req)
	}
	return map[string]interface{}{}, nil
}

func RegOnDeviceControlHandler(fn DeviceControlHandler) {
	gDeviceControlHandler = fn
}

func RegOnDeviceStatusHandler(fn DeviceStatusHandler) {
	gDeviceStatusHandler = fn
}

func RegOnDataTransport(fn TransportDataHandler) {
	gDataTransportHandler = fn
}
