package libgen

import "gitee.com/Puietel/std"

type ModelProperty struct {
	Type     PropertyType `json:"type" validate:"required"`
	Name     string       `json:"name" validate:"required"`
	Restrict Restrict     `json:"restrict" validate:"required"`
}

type ModelInfo PkgInfo

type DeviceModel struct {
	ModelInfo
	ModelProperties []*ModelProperty `json:"properties"`
}

type Device struct {
	ModelInfo  *ModelInfo             `json:"modelInfo" validate:"required"`
	Id         string                 `json:"devId" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	Room       string                 `json:"room" validate:"required"`
	Properties map[string]interface{} `json:"properties"`
	MetaData   map[string]interface{} `json:"metadata"`
}

func NewDeviceModel(pkg string, name string) *DeviceModel {
	return &DeviceModel{
		ModelInfo: ModelInfo{
			Package: pkg,
			Name:    name,
		},
		ModelProperties: make([]*ModelProperty, 0),
	}
}

func (this *DeviceModel) AddModelProperty(tp PropertyType, name string, restrict Restrict) {
	std.Assert(len(name) != 0, "empty name")
	p := &ModelProperty{
		Type:     tp,
		Name:     name,
		Restrict: restrict,
	}
	this.ModelProperties = append(this.ModelProperties, p)
}
