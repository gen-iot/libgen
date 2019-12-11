package libgen

import (
	"errors"
	"github.com/gen-iot/rpcx/v2"
	"github.com/gen-iot/std"
	"log"
	"time"
)

type DeviceCommandHandler func(req *OnDeviceCommandRequest) (std.JsonObject, error)
type DeviceStatusHandler func(notify *DeviceStatusInfo)
type InvokeServiceHandler func(req *OnServiceInvokedRequest) (std.JsonObject, error)
type DeviceIDLENotifyHandler func(req *NotifyDeviceIDLERequest) error

var gDeviceCommandHandler DeviceCommandHandler
var gDeviceStatusHandler DeviceStatusHandler
var gInvokeServiceHandler InvokeServiceHandler
var gDeviceIDLENotifyHandler DeviceIDLENotifyHandler

//noinspection ALL
func pong(ctx rpcx.Context, req *Ping) (*Pong, error) {
	log.Println("receive ping req.time =", req.Time, " delta is ", time.Now().Sub(req.Time))
	return &Pong{Time: time.Now(), Msg: "client pong"}, nil
}

var errAppNotImpControl = errors.New("app not support control device yet")

//noinspection ALL
func onDeviceCommand(ctx rpcx.Context, req *OnDeviceCommandRequest) (std.JsonObject, error) {
	if gDeviceCommandHandler != nil {
		return gDeviceCommandHandler(req)
	}
	return nil, errAppNotImpControl
}

//noinspection ALL
func onDeviceStatusDelivery(ctx rpcx.Context, notify *DeviceStatusInfo) error {
	if gDeviceStatusHandler != nil {
		go gDeviceStatusHandler(notify)
	}
	return nil
}

//noinspection ALL
func onServiceInvoke(ctx rpcx.Context, req *OnServiceInvokedRequest) (std.JsonObject, error) {
	if gInvokeServiceHandler != nil {
		return gInvokeServiceHandler(req)
	}
	return std.NewJsonObject(), nil
}

//noinspection ALL
func onDeviceIDLENotify(ctx rpcx.Context, req *NotifyDeviceIDLERequest) error {
	if gDeviceIDLENotifyHandler != nil {
		return gDeviceIDLENotifyHandler(req)
	}
	return nil
}

//noinspection ALL
func RegOnDeviceCommandHandler(fn DeviceCommandHandler) {
	gDeviceCommandHandler = fn
}

//noinspection ALL
func RegOnDeviceStatusHandler(fn DeviceStatusHandler) {
	gDeviceStatusHandler = fn
}

//noinspection ALL
func RegOnInvokeServiceHandler(fn InvokeServiceHandler) {
	gInvokeServiceHandler = fn
}

//noinspection ALL
func RegOnDeviceIDLENotify(fn DeviceIDLENotifyHandler) {
	gDeviceIDLENotifyHandler = fn
}
