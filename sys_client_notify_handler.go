package libgen

import (
	"errors"
	"github.com/gen-iot/rpcx"
	"log"
	"time"
)

type DeviceControlHandler func(req *ControlDeviceRequest) error
type DeviceStatusHandler func(notify *DeviceStatusInfo)
type TransportDataHandler func(req *TransportDataRequest) (map[string]interface{}, error)

var gDeviceControlHandler DeviceControlHandler
var gDeviceStatusHandler DeviceStatusHandler
var gDataTransportHandler TransportDataHandler

//noinspection ALL
func pong(ctx rpcx.Context, req *Ping) (*Pong, error) {
	log.Println("receive ping req.time =", req.Time, " delta is ", time.Now().Sub(req.Time))
	return &Pong{Time: time.Now(), Msg: "client pong"}, nil
}

var errAppNotImpControl = errors.New("app not support control device yet")

//noinspection ALL
func onDeviceControl(ctx rpcx.Context, req *ControlDeviceRequest) error {
	if gDeviceControlHandler != nil {
		return gDeviceControlHandler(req)
	}
	return errAppNotImpControl
}

//noinspection ALL
func onDeviceStatusDelivery(ctx rpcx.Context, notify *DeviceStatusInfo) error {
	if gDeviceStatusHandler != nil {
		go gDeviceStatusHandler(notify)
	}
	return nil
}

//noinspection ALL
func onDataTransport(ctx rpcx.Context, req *TransportDataRequest) (map[string]interface{}, error) {
	if gDataTransportHandler != nil {
		return gDataTransportHandler(req)
	}
	return map[string]interface{}{}, nil
}

//noinspection ALL
func RegOnDeviceControlHandler(fn DeviceControlHandler) {
	gDeviceControlHandler = fn
}

//noinspection ALL
func RegOnDeviceStatusHandler(fn DeviceStatusHandler) {
	gDeviceStatusHandler = fn
}

//noinspection ALL
func RegOnDataTransport(fn TransportDataHandler) {
	gDataTransportHandler = fn
}
