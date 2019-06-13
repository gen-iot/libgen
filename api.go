package libgen

type GenCommand uint8

const (
	CmdDeclareDeviceModel GenCommand = iota
	CmdRemoveDeviceModels
	CmdUpdateDeviceModel
	CmdRegisterDevices
	CmdRemoveDevices
	CmdUpdateDeviceRoom
	CmdUpdateDeviceName
	CmdUpdateDeviceStatus
	CmdOnDeviceControl
	CmdFetchDevices
	CmdDeviceControl
	CmdOnDeviceStatusChanged
)

type API interface {
	//declare device models,only device model declared can be used in device
	//ps:this function need some privilege//TODO define privilege
	DeclareDeviceModel(models []DeviceModel) error

	//remove device models
	//ps:this function need some privilege//TODO define privilege
	RemoveDeviceModels(modIds []string) error

	//update  device model
	//ps:this function need some privilege//TODO define privilege
	UpdateDeviceModel(modId string, model DeviceModel) error

	//register devices,
	//each device's modelId must be filled with device model declared ahead
	//ps:this function need some privilege//TODO define privilege
	RegisterDevices(devices []Device) error
	//remove devices
	//ps:this function need some privilege//TODO define privilege
	RemoveDevices(devIds []string) error

	//update device room
	//ps:this function need some privilege//TODO define privilege
	UpdateDeviceRoom(devId string, newRoom string) error

	//update device name
	//ps:this function need some privilege//TODO define privilege
	UpdateDeviceName(devId string, newName string) error

	//update device status properties,
	//statusProps must contains all status properties defined in DeviceModel
	//ps:this function need some privilege//TODO define privilege
	UpdateDeviceStatus(devId string, statusProps map[string]interface{}) error

	//device control callback function,
	//this function will be called while someone want control device on your domain
	//ps:this function need some privilege//TODO define privilege
	OnDeviceControl(devId string, cmdProps map[string]interface{}) error

	//fetch all devices
	//ps:this function need some privilege//TODO define privilege
	FetchDevices()

	//control device,control some device in other domain
	//ps:this function need some privilege//TODO define privilege
	DeviceControl(domain string, deviceId string, cmdProps map[string]interface{}) error

	//this is a notify callback function,
	//this function will be called while devices'status which in other domain changed
	//ps:this function need some privilege//TODO define privilege
	OnDeviceStatusChanged(domain string, deviceId string, statusProps map[string]interface{}) error
}
