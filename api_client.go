//+build !server

package libgen

import (
	"gitee.com/gen-iot/rpcx"
	"github.com/pkg/errors"
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

func (this *ApiClientImpl) callWrapper(method string, req interface{}, res interface{}) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call(ApiCallTimeout, method, req, res)
}

func (this *ApiClientImpl) DeclareDeviceModel(req *DeclareDeviceModelRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("DeclareDeviceModel", req, res)
	return res, err
}

func (this *ApiClientImpl) RemoveDeviceModels(req *RemoveDeviceModelsRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("RemoveDeviceModels", req, res)
	return res, err
}

func (this *ApiClientImpl) RegisterDevices(req *RegisterDevicesRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("RegisterDevices", req, res)
	return res, err
}

func (this *ApiClientImpl) RemoveDevices(req *RemoveDevicesRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("RemoveDevices", req, res)
	return res, err
}

func (this *ApiClientImpl) UpdateDeviceInfo(req *UpdateDeviceInfoRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("UpdateDeviceInfo", req, res)
	return res, err
}

func (this *ApiClientImpl) SetDeviceOnline(req *SetOnlineRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("SetDeviceOnline", req, res)
	return res, err
}

func (this *ApiClientImpl) ReportDeviceStatus(req *ReportDeviceStatusRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("ReportDeviceStatus", req, res)
	return res, err
}

func (this *ApiClientImpl) FetchDevices(req *FetchDevicesRequest) (*FetchDevicesResponse, error) {
	res := new(FetchDevicesResponse)
	err := this.callWrapper("FetchDevices", req, res)
	return res, err
}

func (this *ApiClientImpl) ControlDevice(req *ControlDeviceRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("ControlDevice", req, res)
	return res, err
}

func (this *ApiClientImpl) Ping(req *Ping) (*Pong, error) {
	res := new(Pong)
	err := this.callWrapper("Ping", req, res)
	return res, err
}

var emptyReq = make(map[string]interface{})

func (this *ApiClientImpl) SystemSummary() (*SystemSummaryResponse, error) {
	out := new(SystemSummaryResponse)
	err := this.callWrapper("SystemSummary", emptyReq, out)
	return out, err
}

func (this *ApiClientImpl) ListRooms() (*ListRoomsResponse, error) {
	out := new(ListRoomsResponse)
	err := this.callWrapper("ListRooms", emptyReq, out)
	return out, err
}
