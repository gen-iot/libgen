package libgen

import (
	"fmt"
	"time"
)

type PkgInfo struct {
	Package string `json:"package" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

type DeclareDeviceModelRequest struct {
	Model           *DeviceModel `json:"model" validate:"required"`
	OverrideIfExist bool         `json:"overrideIfExist"`
}

type RemoveDeviceModelsRequest struct {
	ModelNames []string `json:"modelNames" validate:"required,gt=0"`
}

type RegisterDevicesRequest struct {
	Devices []*Device `json:"devices" validate:"required,gt=0,dive,required"`
}

type RemoveDevicesRequest struct {
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

type ControlDeviceRequest struct {
	PkgInfo    PkgInfo                `json:"pkgInfo" validate:"required"`
	Id         string                 `json:"id" validate:"required"`
	CtrlFunc   string                 `json:"ctrlFunc"`
	CtrlParams map[string]interface{} `json:"ctrlParams" validate:"required,gt=0"`
}

type OnDeviceControlRequest struct {
	Id         string                 `json:"id" validate:"required"`
	CtrlParams map[string]interface{} `json:"ctrlParams" validate:"required,gt=0"`
}

type DeviceStatusInfo struct {
	PkgInfo   PkgInfo                `json:"pkgInfo"`
	ModelInfo *ModelInfo             `json:"modelInfo"`
	Id        string                 `json:"id" validate:"required"`
	Name      string                 `json:"name" validate:"required"`
	Room      string                 `json:"room" validate:"required"`
	Status    map[string]interface{} `json:"status"`
	Online    bool                   `json:"online"`
}

type SetOnlineRequest struct {
	DeviceIds []string `json:"deviceIds" validate:"required"`
	Online    bool     `json:"online"`
}

type HandshakeRequest struct {
	PkgInfo
	AccessToken string `json:"accessToken"`
}

type Ping struct {
	Time time.Time `json:"time"`
	Msg  string    `json:"msg"`
}

type TransportDataRequest struct {
	Sender *PkgInfo `json:"sender" validate:"required"`
	Data   []byte   `json:"data" validate:"required"`
}

func (this *Ping) String() string {
	return fmt.Sprintf("time = %s , extraMsg = %s", this.Time.String(), this.Msg)
}

type Pong = Ping

type SystemSummaryResponse struct {
	GenId   string `json:"genId"`
	Version string `json:"version"`
}

type ListRoomsResponse struct {
	Rooms []string `json:"rooms"`
}

type DeviceModelRuntime struct {
	ModelInfo
	StatusProperties  map[string]*StatusProperty  `json:"statusProperties"`
	CommandProperties map[string]*ControlProperty `json:"commandProperties"`
}

type ListDeviceModelRequest struct {
	Includes []*ModelInfo `json:"includes" validate:"omitempty,gt=0,dive,required"`
}
type ListDeviceModelResponse struct {
	Models []*DeviceModelRuntime `json:"models"`
}

type ListDevicesByRoomRequest struct {
	Rooms    []string     `json:"rooms" validate:"gt=0"`
	Includes []*ModelInfo `json:"includes" validate:"omitempty,gt=0,dive,required"`
}

type RoomDeviceResultItem struct {
	Room        string              `json:"room"`
	RoomDevices []*DeviceStatusInfo `json:"devices"`
	Error       string              `json:"error,omitempty"`
}

type ListDevicesByRoomResult struct {
	RoomDevices []*RoomDeviceResultItem `json:"roomDevices" validate:"omitempty,gt=0,dive,required"`
}

// client side
type RpcApiClient interface {
	//ping
	Ping(req *Ping) (*Pong, error)

	//declare device models,only device model declared can be used in device
	DeclareDeviceModel(req *DeclareDeviceModelRequest) error

	//remove device models
	RemoveDeviceModels(req *RemoveDeviceModelsRequest) error

	//register devices,
	//each device's modelId must be filled with device model declared ahead
	RegisterDevices(req *RegisterDevicesRequest) error

	//remove devices
	RemoveDevices(req *RemoveDevicesRequest) error

	//remove all app devices
	RemoveAppDevice() error

	//set online
	SetDeviceOnline(req *SetOnlineRequest) error

	//update device info,
	UpdateDeviceInfo(req *UpdateDeviceInfoRequest) error

	//update device status,
	ReportDeviceStatus(req *ReportDeviceStatusRequest) error

	//control devices
	ControlDevice(req *ControlDeviceRequest) error

	SystemSummary() (*SystemSummaryResponse, error)

	ListRooms() (*ListRoomsResponse, error)

	ListDeviceModels(req *ListDeviceModelRequest) (*ListDeviceModelResponse, error)

	ListDevicesByRoom(req *ListDevicesByRoomRequest) (*ListDevicesByRoomResult, error)
}
