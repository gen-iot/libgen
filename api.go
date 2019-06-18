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
	Models []*DeviceModel `json:"models" validate:"required,gt=0"`
}

type RemoveDeviceModelsRequest struct {
	BaseRequest
	Ids []string `json:"ids" validate:"required,gt=0"`
}

type UpdateDeviceModelRequest struct {
	BaseRequest
	Id    string       `json:"id" validate:"required"`
	Model *DeviceModel `json:"model" validate:"required"`
}

type RegisterDevicesRequest struct {
	BaseRequest
	Devices []*Device `json:"devices" validate:"required,gt=0"`
}

type RemoveDevicesRequest struct {
	BaseRequest
	Ids []string `json:"ids" validate:"required,gt=0"`
}

type UpdateDeviceRequest struct {
	Id         string                 `json:"id" validate:"required"`
	Name       *string                `json:"name"`
	Room       *string                `json:"room"`
	Properties map[string]interface{} `json:"properties"`
	MetaData   map[string]interface{} `json:"metaData"`
}

type FetchDevicesRequest struct {
	BaseRequest
	//if id is not nil or empty will be as the only query condition
	Id *string `json:"id"`
	//filter condition , if filed below not nil or empty will be as '&&' query condition
	Name   *string `json:"name"`
	Room   *string `json:"room"`
	Domain *string `json:"domain"`
}

type FetchDevicesResponse struct {
	BaseResponse
	Devices []*DeviceWithDomain
}

type ControlDeviceRequest struct {
	BaseRequest
	Domain     string                 `json:"domain" validate:"required"`
	Id         string                 `json:"id" validate:"required"`
	CtrlParams map[string]interface{} `json:"ctrlParams" validate:"required"`
}

type ControlDeviceResponse struct {
	BaseResponse
}

type Ping struct {
	Time time.Time `json:"time"`
	Msg  string    `json:"msg"`
}

func (this *Ping) String() string {
	return fmt.Sprintf("time = %s , extraMsg = %s", this.Time.String(), this.Msg)
}

type Pong = Ping

type RpcApi interface {
	//declare device models,only device model declared can be used in device
	DeclareDeviceModel(req *DeclareDeviceModelRequest) (*BaseResponse, error)

	//remove device models
	RemoveDeviceModels(req *RemoveDeviceModelsRequest) (*BaseResponse, error)

	//update  device model
	UpdateDeviceModel(req *UpdateDeviceModelRequest) (*BaseResponse, error)

	//register devices,
	//each device's modelId must be filled with device model declared ahead
	RegisterDevices(req *RegisterDevicesRequest) (*BaseResponse, error)

	//remove devices
	RemoveDevices(req *RemoveDevicesRequest) (*BaseResponse, error)

	//update device,
	UpdateDevice(req *UpdateDeviceRequest) (*BaseResponse, error)

	//fetch devices
	FetchDevices(req *FetchDevicesRequest) (*FetchDevicesResponse, error)

	//control devices
	DeviceControl(req *ControlDeviceRequest) (*BaseResponse, error)

	//ping
	Ping(req *Ping) (*Pong, error)
}
