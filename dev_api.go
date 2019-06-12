package libgen

//declare device models,only device model declared can be used in device
//ps:this function need some privilege//TODO define privilege
func DeclareDeviceModel(models ...DeviceModel) error {

	return nil
}

//remove device models
//ps:this function need some privilege//TODO define privilege
func RemoveDeviceModels(modIds ...string) error {

	return nil
}

//update  device model
//ps:this function need some privilege//TODO define privilege
func UpdateDeviceModel(modId string, model DeviceModel) error {

	return nil
}

//register devices,
//each device's modelId must be filled with device model declared ahead
//ps:this function need some privilege//TODO define privilege
func RegisterDevices(devices ...Device) error {
	return nil
}

//remove devices
//ps:this function need some privilege//TODO define privilege
func RemoveDevices(devIds ...string) error {
	return nil
}

//update device room
//ps:this function need some privilege//TODO define privilege
func UpdateDeviceRoom(devId string, newRoom string) error {
	return nil
}

//update device name
//ps:this function need some privilege//TODO define privilege
func UpdateDeviceName(devId string, newName string) error {
	return nil
}

//update device status properties,
//statusProps must contains all status properties defined in DeviceModel
//ps:this function need some privilege//TODO define privilege
func UpdateDeviceStatus(devId string, statusProps map[string]interface{}) error {
	return nil
}

//device control callback function,
//this function will be called while someone want control device on your domain
//ps:this function need some privilege//TODO define privilege
func OnDeviceControl(devId string, cmdProps map[string]interface{}) error {

	return nil
}

//fetch all devices
//ps:this function need some privilege//TODO define privilege
func FetchDevices() {

}

//control device,control some device in other domain
//ps:this function need some privilege//TODO define privilege
func DeviceControl(domain string, deviceId string, cmdProps map[string]interface{}) error {

	return nil
}

//this is a notify callback function,
//this function will be called while devices'status which in other domain changed
//ps:this function need some privilege//TODO define privilege
func OnDeviceStatusChanged(domain string, deviceId string, statusProps map[string]interface{}) error {

	return nil
}
