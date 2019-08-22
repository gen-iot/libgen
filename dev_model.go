package libgen

import "github.com/gen-iot/std"

type ModelInfo = PkgInfo

type DeviceModel struct {
	ModelInfo
	StatusProperties  []*StatusProperty  `json:"statusProperties"`
	ControlProperties []*ControlProperty `json:"controlGroups"`
}

type Device struct {
	ModelInfo  *ModelInfo             `json:"modelInfo" validate:"required"`
	Id         string                 `json:"devId" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	Room       string                 `json:"room" validate:"required"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	MetaData   map[string]interface{} `json:"metadata,omitempty"`
}

//noinspection ALL
func NewDeviceModel(pkg string, name string) *DeviceModel {
	return &DeviceModel{
		ModelInfo: ModelInfo{
			Package: pkg,
			Name:    name,
		},
		StatusProperties:  make([]*StatusProperty, 0),
		ControlProperties: make([]*ControlProperty, 0),
	}
}

func (this *DeviceModel) AddModelProperty(name string, restrict Restrict) {
	std.Assert(len(name) != 0, "empty name")
	p := &StatusProperty{
		Name:     name,
		Restrict: restrict,
	}
	this.StatusProperties = append(this.StatusProperties, p)
}
