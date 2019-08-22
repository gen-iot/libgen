//+build !server

package libgen

import (
	"errors"
	"github.com/gen-iot/rpcx"
	"sync"
)

type ApiClientImpl struct {
	rwLock   *sync.RWMutex
	callable rpcx.Callable
}

func NewApiClientImpl() *ApiClientImpl {
	return &ApiClientImpl{
		rwLock:   &sync.RWMutex{},
		callable: nil,
	}
}

func (this *ApiClientImpl) setCallable(callable rpcx.Callable) {
	this.rwLock.Lock()
	defer this.rwLock.Unlock()
	this.callable = callable
}

func (this *ApiClientImpl) getCallable() rpcx.Callable {
	this.rwLock.RLock()
	defer this.rwLock.RUnlock()
	return this.callable
}

var errConnectionClosed = errors.New("the connection to gen not established")

func (this *ApiClientImpl) callWrapper(method string) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call(ApiCallTimeout, method)
}

func (this *ApiClientImpl) call1Wrapper(method string, req interface{}) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call1(ApiCallTimeout, method, req)
}

func (this *ApiClientImpl) call5Wrapper(method string, req interface{}, res interface{}) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call5(ApiCallTimeout, method, req, res)
}

func (this *ApiClientImpl) call3Wrapper(method string, res interface{}) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call3(ApiCallTimeout, method, res)
}

func (this *ApiClientImpl) SystemSummary() (*SystemSummaryResponse, error) {
	out := new(SystemSummaryResponse)
	err := this.call3Wrapper("SystemSummary", out)
	return out, err
}

func (this *ApiClientImpl) Ping(req *Ping) (*Pong, error) {
	res := new(Pong)
	err := this.call5Wrapper("Ping", req, res)
	return res, err
}

func (this *ApiClientImpl) DeclareDeviceModel(req *DeclareDeviceModelRequest) error {
	return this.call1Wrapper("DeclareDeviceModel", req)
}

func (this *ApiClientImpl) RemoveDeviceModels(req *RemoveDeviceModelsRequest) error {
	return this.call1Wrapper("RemoveDeviceModels", req)
}

func (this *ApiClientImpl) RegisterDevices(req *RegisterDevicesRequest) error {
	return this.call1Wrapper("RegisterDevices", req)
}

func (this *ApiClientImpl) RemoveDevices(req *RemoveDevicesRequest) error {
	return this.call1Wrapper("RemoveDevices", req)
}

func (this *ApiClientImpl) RemoveAppDevice() error {
	return this.callWrapper("RemoveAppDevice")
}

func (this *ApiClientImpl) UpdateDeviceInfo(req *UpdateDeviceInfoRequest) error {
	return this.call1Wrapper("UpdateDeviceInfo", req)
}

func (this *ApiClientImpl) SetDeviceOnline(req *SetOnlineRequest) error {
	return this.call1Wrapper("SetDeviceOnline", req)
}

func (this *ApiClientImpl) ReportDeviceStatus(req *ReportDeviceStatusRequest) error {
	return this.call1Wrapper("ReportDeviceStatus", req)
}

func (this *ApiClientImpl) ControlDevice(req *CommandDeviceRequest) error {
	return this.call1Wrapper("ControlDevice", req)
}

func (this *ApiClientImpl) ListRooms() (*ListRoomsResponse, error) {
	out := new(ListRoomsResponse)
	err := this.call3Wrapper("ListRooms", out)
	return out, err
}

func (this *ApiClientImpl) ListDeviceModels(req *ListDeviceModelRequest) (*ListDeviceModelResponse, error) {
	out := new(ListDeviceModelResponse)
	err := this.call5Wrapper("ListDeviceModels", req, out)
	return out, err
}

func (this *ApiClientImpl) ListDevicesByRoom(req *ListDevicesByRoomRequest) (*ListDevicesByRoomResult, error) {
	res := new(ListDevicesByRoomResult)
	err := this.call5Wrapper("ListDevicesByRoom", req, res)
	return res, err
}
