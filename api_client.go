package libgen

type ApiClientImpl struct {
}

func (this *ApiClientImpl) DeclareDeviceModel(req *DeclareDeviceModelRequest) (*BaseResponse, error) {
	res := &BaseResponse{}
	err := gCallable.Call(ApiCallTimeout, "DeclareDeviceModel", req, res)
	return res, err
}

func (this *ApiClientImpl) RemoveDeviceModels(req *RemoveDeviceModelsRequest) (*BaseResponse, error) {

	panic("implement me")
}

func (this *ApiClientImpl) UpdateDeviceModel(req *UpdateDeviceModelRequest) (*BaseResponse, error) {

	panic("implement me")
}

func (this *ApiClientImpl) RegisterDevices(req *RegisterDevicesRequest) (*BaseResponse, error) {

	panic("implement me")
}

func (this *ApiClientImpl) RemoveDevices(req *RemoveDevicesRequest) (*BaseResponse, error) {

	panic("implement me")
}

func (this *ApiClientImpl) UpdateDevice(req *UpdateDeviceRequest) (*BaseResponse, error) {

	panic("implement me")
}

func (this *ApiClientImpl) FetchDevices(req *FetchDevicesRequest) (*FetchDevicesResponse, error) {

	panic("implement me")
}

func (this *ApiClientImpl) DeviceControl(req *ControlDeviceRequest) (*BaseResponse, error) {

	panic("implement me")
}

func (this *ApiClientImpl) Ping(req *Ping) (*Pong, error) {
	res := new(Pong)
	err := gCallable.Call(ApiCallTimeout, "Ping", req, res)
	return res, err
}
