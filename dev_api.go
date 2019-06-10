package libgen

import "fmt"

func Init() {
	fmt.Println("LIBGEN INIT")
}

func AddDevice() error {
	return nil
}

func RemoveDevices(devIds ...string) error {
	return nil
}

func UpdateDeviceName(devId string, newName string) error {
	return nil
}

