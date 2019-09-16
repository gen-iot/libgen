package libgen

import (
	"fmt"
	"github.com/gen-iot/std"
	"time"
)

//noinspection ALL
func NewDeviceModel(pkg string, name string) *DeviceModel {
	return &DeviceModel{
		ModelInfo: ModelInfo{
			Package: pkg,
			Name:    name,
		},
		StatusProperties:  make([]Restrict, 0),
		CommandProperties: make([]*CommandProperty, 0),
	}
}
func (this *DeviceModel) AddModelProperty(restrict Restrict) {
	this.StatusProperties = append(this.StatusProperties, restrict)
}

type PkgInfo struct {
	Package string `json:"package" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

type ModelInfo = PkgInfo

type DeviceModel struct {
	ModelInfo
	StatusProperties  []Restrict         `json:"statusProperties" validate:"omitempty,dive,required"`
	CommandProperties []*CommandProperty `json:"commandProperties" validate:"omitempty,dive,required"`
}
type RestrictPersist std.JsonObject

type CommandPropertyPersist struct {
	FuncName       string            `json:"command" validate:"required"`
	ParamRestricts []RestrictPersist `json:"restricts" validate:"omitempty,dive,required"`
}
type DeviceModelPersist struct {
	ModelInfo
	StatusProperties  []RestrictPersist         `json:"statusProperties"`
	CommandProperties []*CommandPropertyPersist `json:"commandProperties"`
}
type Device struct {
	ModelInfo *ModelInfo     `json:"modelInfo" validate:"required"`
	Id        string         `json:"id" validate:"required"`
	Name      string         `json:"name" validate:"required"`
	Room      string         `json:"room" validate:"required"`
	Category  string         `json:"category"`
	MetaData  std.JsonObject `json:"metadata,omitempty"`
}

type DeviceStatusInfo struct {
	*Device
	PkgInfo *PkgInfo       `json:"pkgInfo"`
	Online  bool           `json:"online"`
	Status  std.JsonObject `json:"status,omitempty"`
}

type Ping struct {
	Time time.Time `json:"time"`
	Msg  string    `json:"msg"`
}

func (this *Ping) String() string {
	return fmt.Sprintf("time = %s , msg = %s", this.Time.String(), this.Msg)
}

type Pong = Ping

type SystemSummaryResponse struct {
	GenId   string `json:"genId"`
	Version string `json:"version"`
}

type DeclareDeviceModelRequest struct {
	ModelName         string             `json:"modelName" validate:"required"`
	StatusProperties  []Restrict         `json:"statusProperties" validate:"omitempty,dive,required"`
	CommandProperties []*CommandProperty `json:"commandProperties" validate:"omitempty,dive,required"`
	OverrideIfExist   bool               `json:"overrideIfExist"`
}

type RemoveDeviceModelsRequest struct {
	ModelNames []string `json:"modelNames" validate:"required,gt=0,dive,required"`
}
type ListDeviceModelRequest struct {
	Selector []*ModelInfo `json:"filter" validate:"omitempty,dive,required"`
}

type ListDeviceModelResponse struct {
	Models []*DeviceModelPersist `json:"models"`
}

type RegisterDevicesRequest struct {
	Devices []*Device `json:"devices" validate:"required,gt=0,dive,required"`
}

type RemoveDevicesRequest struct {
	Ids []string `json:"ids" validate:"required,gt=0,dive,required"`
}

type UpdateDeviceInfoRequest struct {
	Id       string         `json:"id" validate:"required"`
	Name     *string        `json:"name"`
	Room     *string        `json:"room"`
	MetaData std.JsonObject `json:"metaData"`
}

type ReportDeviceStatusRequest struct {
	Id     string         `json:"id" validate:"required"`
	Status std.JsonObject `json:"status"`
	Online bool           `json:"online"`
}

type SetOnlineRequest struct {
	Ids    []string `json:"ids" validate:"required,gt=0,dive,required"`
	Online bool     `json:"online"`
}

type ListRoomsResponse struct {
	Rooms []string `json:"rooms"`
}

type ListDevicesByRoomRequest struct {
	Rooms          []string     `json:"rooms" validate:"omitempty,dive,required"`
	Filter         []*ModelInfo `json:"filter" validate:"omitempty,dive,required"`
	CategoryFilter []string     `json:"categoryFilter" validate:"omitempty,dive,required"`
}
type RoomDevicesBucket struct {
	Room    string              `json:"room"`
	Devices []*DeviceStatusInfo `json:"devices"`
}
type ListDevicesByRoomResult struct {
	RoomDevices []*RoomDevicesBucket `json:"rooms" validate:"omitempty,gt=0,dive,required"`
}

type FindDeviceByIdRequest struct {
	PkgInfo  PkgInfo `json:"pkgInfo" validate:"required"`
	DeviceId string  `json:"deviceId" validate:"dive,required"`
}

type FindDeviceByIdResponse = DeviceStatusInfo

type CommandDeviceRequest struct {
	PkgInfo PkgInfo        `json:"pkgInfo" validate:"required"`
	Id      string         `json:"id" validate:"required"`
	Command string         `json:"command"`
	Params  std.JsonObject `json:"params"`
}

type OnDeviceCommandRequest struct {
	Id      string         `json:"id" validate:"required"`
	Command string         `json:"command" validate:"required"`
	Params  std.JsonObject `json:"params"`
}

type HandshakeRequest struct {
	PkgInfo
	ApiAccessToken string `json:"apiAccessToken"`
}

type TransportDataRequest struct {
	Sender *PkgInfo `json:"sender" validate:"required"`
	Data   []byte   `json:"data" validate:"required"`
}

type NotifyDeviceIDLERequest struct {
	Ids []string `json:"Ids"`
}
