package libgen

import (
	"gitee.com/Puietel/std"
	"log"
	"testing"
	"time"
)

func TestDeclareDeviceModel(t *testing.T) {
	models := make([]*DeviceModel, 0)
	models = append(models, &DeviceModel{
		Id:   "0x41",
		Name: "普通灯",
		ModelProperties: []*ModelProperty{
			{
				Type:     Command,
				Name:     "power",
				Restrict: NewI32EnumLimiter([]int32{1, 2}),
			},
			{
				Type:     Status,
				Name:     "power",
				Restrict: NewI32EnumLimiter([]int32{1, 2}),
			},
		},
	})
	req := &DeclareDeviceModelRequest{
		Models: models,
	}
	rsp, err := req.Wait(time.Second * 10)
	std.AssertError(err, "Wait Response")
	if !rsp.Ok() {
		log.Println("declare models error :", rsp.Error())
	}
	log.Println("declare models ok")
}
