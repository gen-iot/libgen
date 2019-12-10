//+build !server

package libgen

import (
	"errors"
	"github.com/gen-iot/rpcx/v2"
	"github.com/gen-iot/std"
	"sync"
)

type LinkMethod string

const (
	Handshake   LinkMethod = "Handshake"
	DebugAttach LinkMethod = "DebugAttach"
)

const (
	// supported func list
	kDeliveryDeviceStatus = "DeliveryDeviceStatus"
	kPing                 = "Ping"
	kInvokeService        = "InvokeService"
	kNotifyDeviceIDLE     = "NotifyDeviceIDLE"
	kSystemSummary        = "SystemSummary"
	kDeclareDeviceModel   = "DeclareDeviceModel"
	kRemoveDeviceModels   = "RemoveDeviceModels"
	kListDeviceModels     = "ListDeviceModels"
	kRegisterDevices      = "RegisterDevices"
	kRemoveDevices        = "RemoveDevices"
	kRemoveAppDevice      = "RemoveAppDevice"
	kUpdateDeviceInfo     = "UpdateDeviceInfo"
	kReportDeviceStatus   = "ReportDeviceStatus"
	kSetDeviceOnline      = "SetDeviceOnline"
	kListRooms            = "ListRooms"
	kListDevicesByRoom    = "ListDevicesByRoom"
	kFindDeviceById       = "FindDeviceById"
	kCommandDevice        = "CommandDevice"
)

func NewApiClientImpl() RpcApiClient {
	return &apiClientImpl{
		rwLock:   &sync.RWMutex{},
		callable: nil,
	}
}

type apiClientImpl struct {
	rwLock   *sync.RWMutex
	callable rpcx.Callable
}

func (this *apiClientImpl) setCallable(callable rpcx.Callable) {
	this.rwLock.Lock()
	defer this.rwLock.Unlock()
	this.callable = callable
}
func (this *apiClientImpl) getCallable() rpcx.Callable {
	this.rwLock.RLock()
	defer this.rwLock.RUnlock()
	return this.callable
}

var errConnectionClosed = errors.New("the connection to gen not established")

func (this *apiClientImpl) callWrapper(method string) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call(ApiCallTimeout, method)
}

func (this *apiClientImpl) call1Wrapper(method string, req interface{}) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call1(ApiCallTimeout, method, req)
}

func (this *apiClientImpl) call5Wrapper(method string, req interface{}, res interface{}) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call5(ApiCallTimeout, method, req, res)
}

func (this *apiClientImpl) call3Wrapper(method string, res interface{}) error {
	callable := this.getCallable()
	if callable == nil {
		return errConnectionClosed
	}
	return callable.Call3(ApiCallTimeout, method, res)
}

func (this *apiClientImpl) Ping(req *Ping) (*Pong, error) {
	res := new(Pong)
	err := this.call5Wrapper(kPing, req, res)
	return res, err
}

func (this *apiClientImpl) SystemSummary() (*SystemSummaryResponse, error) {
	out := new(SystemSummaryResponse)
	err := this.call3Wrapper(kSystemSummary, out)
	return out, err
}

func (this *apiClientImpl) DeclareDeviceModel(req *DeclareDeviceModelRequest) error {
	return this.call1Wrapper(kDeclareDeviceModel, req)
}

func (this *apiClientImpl) RemoveDeviceModels(req *RemoveDeviceModelsRequest) error {
	return this.call1Wrapper(kRemoveDeviceModels, req)
}

func (this *apiClientImpl) ListDeviceModels(req *ListDeviceModelRequest) (*ListDeviceModelResponse, error) {
	out := new(ListDeviceModelResponse)
	err := this.call5Wrapper(kListDeviceModels, req, out)
	return out, err
}

func (this *apiClientImpl) RegisterDevices(req *RegisterDevicesRequest) error {
	return this.call1Wrapper(kRegisterDevices, req)
}

func (this *apiClientImpl) RemoveDevices(req *RemoveDevicesRequest) error {
	return this.call1Wrapper(kRemoveDevices, req)
}

func (this *apiClientImpl) RemoveAppDevice() error {
	return this.callWrapper(kRemoveAppDevice)
}

func (this *apiClientImpl) UpdateDeviceInfo(req *UpdateDeviceInfoRequest) error {
	return this.call1Wrapper(kUpdateDeviceInfo, req)
}

func (this *apiClientImpl) ReportDeviceStatus(req *ReportDeviceStatusRequest) error {
	return this.call1Wrapper(kReportDeviceStatus, req)
}

func (this *apiClientImpl) SetDeviceOnline(req *SetOnlineRequest) error {
	return this.call1Wrapper(kSetDeviceOnline, req)
}

func (this *apiClientImpl) ListRooms() (*ListRoomsResponse, error) {
	out := new(ListRoomsResponse)
	err := this.call3Wrapper(kListRooms, out)
	return out, err
}
func (this *apiClientImpl) ListDevicesByRoom(req *ListDevicesByRoomRequest) (*ListDevicesByRoomResult, error) {
	res := new(ListDevicesByRoomResult)
	err := this.call5Wrapper(kListDevicesByRoom, req, res)
	return res, err
}

func (this *apiClientImpl) FindDeviceById(req *FindDeviceByIdRequest) (*FindDeviceByIdResponse, error) {
	res := new(FindDeviceByIdResponse)
	err := this.call5Wrapper(kFindDeviceById, req, res)
	return res, err
}

func (this *apiClientImpl) CommandDevice(req *CommandDeviceRequest) (std.JsonObject, error) {
	out := std.NewJsonObject()
	err := this.call5Wrapper(kCommandDevice, req, &out)
	return out, err
}
