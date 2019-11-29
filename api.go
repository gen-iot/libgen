package libgen

import (
	"github.com/gen-iot/rpcx/v2"
	"github.com/gen-iot/std"
)

// client side
type RpcApiClient interface {
	setCallable(callable rpcx.Callable)
	getCallable() rpcx.Callable
	//ping
	Ping(req *Ping) (*Pong, error)
	//get gen summary info
	SystemSummary() (*SystemSummaryResponse, error)
	//declare device models,only device model declared can be used in device
	DeclareDeviceModel(req *DeclareDeviceModelRequest) error
	//remove device models
	RemoveDeviceModels(req *RemoveDeviceModelsRequest) error
	//list device models in gen instance
	ListDeviceModels(req *ListDeviceModelRequest) (*ListDeviceModelResponse, error)
	//register devices,
	//each device's modelId must be filled with device model declared ahead
	RegisterDevices(req *RegisterDevicesRequest) error
	//remove devices
	RemoveDevices(req *RemoveDevicesRequest) error
	//remove all app devices
	RemoveAppDevice() error
	//update device info,
	UpdateDeviceInfo(req *UpdateDeviceInfoRequest) error
	//update device status,
	ReportDeviceStatus(req *ReportDeviceStatusRequest) error
	//set devices online status
	SetDeviceOnline(req *SetOnlineRequest) error
	//list all rooms in gen instance
	ListRooms() (*ListRoomsResponse, error)
	//list devices in specific rooms
	ListDevicesByRoom(req *ListDevicesByRoomRequest) (*ListDevicesByRoomResult, error)
	//find device with specific id
	FindDeviceById(req *FindDeviceByIdRequest) (*FindDeviceByIdResponse, error)
	//command devices
	CommandDevice(req *CommandDeviceRequest) (std.JsonObject, error)
}
