package libgen

import "fmt"

func Init() {
	fmt.Println("LIBGEN INIT")
}

//declare device models,only device model declared can be used in device
func DeclareDeviceModel(models ...DeviceModel) error {

	return nil
}

//remove device models
func RemoveDeviceModels(modIds ...string) error {

	return nil
}

//update  device model
func UpdateDeviceModel(modId string, model DeviceModel) error {

	return nil
}

//register devices,
//each device's modelId must be filled with device model declared ahead
func RegisterDevices(devices ...Device) error {
	return nil
}

//remove devices
func RemoveDevices(devIds ...string) error {
	return nil
}

//update device room
func UpdateDeviceRoom(devId string, newRoom string) error {
	return nil
}

//update device name
func UpdateDeviceName(devId string, newName string) error {
	return nil
}

//update device status properties,
//statusProps must contains all status properties defined in DeviceModel
func UpdateDeviceStatus(devId string, statusProps map[string]interface{}) error {
	return nil
}
