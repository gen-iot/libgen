//+build !server

package libgen

type ApiClientImpl struct {
}

func (this *ApiClientImpl) callWrapper(method string, req interface{}, res interface{}) error {
	return getCallable().Call(ApiCallTimeout, method, req, res)
}

func (this *ApiClientImpl) DeclareDeviceModel(req *DeclareDeviceModelRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("DeclareDeviceModel", req, res)
	return res, err
}

func (this *ApiClientImpl) RemoveDeviceModels(req *RemoveDeviceModelsRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("RemoveDeviceModels", req, res)
	return res, err
}

func (this *ApiClientImpl) RegisterDevices(req *RegisterDevicesRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("RegisterDevices", req, res)
	return res, err
}

func (this *ApiClientImpl) RemoveDevices(req *RemoveDevicesRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("RemoveDevices", req, res)
	return res, err
}

func (this *ApiClientImpl) UpdateDeviceInfo(req *UpdateDeviceInfoRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("UpdateDeviceInfo", req, res)
	return res, err
}
func (this *ApiClientImpl) ReportDeviceStatus(req *ReportDeviceStatusRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("ReportDeviceStatus", req, res)
	return res, err
}

func (this *ApiClientImpl) FetchDevices(req *FetchDevicesRequest) (*FetchDevicesResponse, error) {
	res := new(FetchDevicesResponse)
	err := this.callWrapper("FetchDevices", req, res)
	return res, err
}

func (this *ApiClientImpl) DeviceControl(req *ControlDeviceRequest) (*BaseResponse, error) {
	res := new(BaseResponse)
	err := this.callWrapper("DeviceControl", req, res)
	return res, err
}

func (this *ApiClientImpl) Ping(req *Ping) (*Pong, error) {
	res := new(Pong)
	err := gCallable.Call(ApiCallTimeout, "Ping", req, res)
	return res, err
}
