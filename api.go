package libgen

import (
	"fmt"
	"time"
)

type BaseRequest struct {
}

type BaseResponse struct {
}

type DeclareDeviceModelRequest struct {
	BaseRequest
	Model           *DeviceModel `json:"model" validate:"required"`
	OverrideIfExist bool         `json:"overrideIfExist"`
}

type RemoveDeviceModelsRequest struct {
	BaseRequest
	ModelNames []string `json:"modelNames" validate:"required,gt=0"`
}

type RegisterDevicesRequest struct {
	BaseRequest
	Devices []*Device `json:"devices" validate:"required,gt=0"`
}

type RemoveDevicesRequest struct {
	BaseRequest
	Ids []string `json:"ids" validate:"required,gt=0"`
}

type UpdateDeviceInfoRequest struct {
	Id       string                 `json:"id" validate:"required"`
	Name     *string                `json:"name"`
	Room     *string                `json:"room"`
	MetaData map[string]interface{} `json:"metaData"`
}

type ReportDeviceStatusRequest struct {
	Id         string                 `json:"id" validate:"required"`
	Properties map[string]interface{} `json:"properties"`
	Online     bool                   `json:"online"`
}

type FetchDevicesRequest struct {
	BaseRequest
	PkgInfo *PkgInfo `json:"pkgInfo" validate:"omitempty"` // filter condition , if filed below not nil or empty will be as '&&' query condition
	Id      *string  `json:"id"`                           // if id is not nil or empty will be as the only query condition
	Room    *string  `json:"room"`
}

type FetchDevicesResponse struct {
	BaseResponse
	Package string    `json:"package"`
	Devices []*Device `json:"devices"`
}

type PkgInfo struct {
	Package string `json:"package" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

type HandshakeRequest struct {
	BaseRequest
	PkgInfo
	AccessToken string `json:"accessToken"`
}

type ControlDeviceRequest struct {
	BaseRequest
	PkgInfo    PkgInfo                `json:"pkgInfo" validate:"required"`
	Id         string                 `json:"id" validate:"required"`
	CtrlParams map[string]interface{} `json:"ctrlParams" validate:"gt=0"`
}

type OnDeviceControlRequest struct {
	BaseRequest
	Id         string                 `json:"id" validate:"required"`
	CtrlParams map[string]interface{} `json:"ctrlParams" validate:"gt=0"`
}

type DeviceStatusNotify struct {
	PkgInfo      PkgInfo                `json:"pkgInfo"`
	Id           string                 `json:"id" validate:"required"`
	StatusParams map[string]interface{} `json:"ctrlParams" validate:"required"`
}

type ControlDeviceResponse struct {
	BaseResponse
}

type SetOnlineRequest struct {
	BaseRequest
	DeviceIds []string `json:"deviceIds" validate:"required"`
	Online    bool     `json:"online"`
}

type Ping struct {
	Time time.Time `json:"time"`
	Msg  string    `json:"msg"`
}

func (this *Ping) String() string {
	return fmt.Sprintf("time = %s , extraMsg = %s", this.Time.String(), this.Msg)
}

type Pong = Ping

// client side
type RpcApiClient interface {
	//declare device models,only device model declared can be used in device
	DeclareDeviceModel(req *DeclareDeviceModelRequest) (*BaseResponse, error)

	//remove device models
	RemoveDeviceModels(req *RemoveDeviceModelsRequest) (*BaseResponse, error)

	//register devices,
	//each device's modelId must be filled with device model declared ahead
	RegisterDevices(req *RegisterDevicesRequest) (*BaseResponse, error)

	//remove devices
	RemoveDevices(req *RemoveDevicesRequest) (*BaseResponse, error)

	//set online
	SetOnline(req *SetOnlineRequest) (*BaseResponse, error)

	//update device info,
	UpdateDeviceInfo(req *UpdateDeviceInfoRequest) (*BaseResponse, error)

	//update device status,
	ReportDeviceStatus(req *ReportDeviceStatusRequest) (*BaseResponse, error)

	//fetch devices
	FetchDevices(req *FetchDevicesRequest) (*FetchDevicesResponse, error)

	//control devices
	ControlDevice(req *ControlDeviceRequest) (*BaseResponse, error)

	//ping
	Ping(req *Ping) (*Pong, error)
}
