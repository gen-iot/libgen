package libgen

import "time"

type DeclareDeviceModelRequest struct {
	*BaseRequest
	Models []*DeviceModel `json:"models"`
}

type DeclareDeviceModelResponse struct {
	*BaseResponse
}

func (this *DeclareDeviceModelRequest) Data() interface{} {

	return this
}

func NewDeclareDeviceModelRequest() *DeclareDeviceModelRequest {
	return &DeclareDeviceModelRequest{
		BaseRequest: NewBaseRequest(CmdDeclareDeviceModel),
	}
}

func (this *DeclareDeviceModelRequest) Wait(timeout time.Duration) (rsp *DeclareDeviceModelResponse, err error) {
	rsp = new(DeclareDeviceModelResponse)
	err = this.getResponse(timeout, this, rsp)
	return
}
