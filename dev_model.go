package libgen

type ModelInfo = PkgInfo

type DeviceModel struct {
	ModelInfo
	StatusProperties  []Restrict         `json:"statusProperties" validate:"omitempty,dive,required"`
	CommandProperties []*CommandProperty `json:"commandProperties" validate:"omitempty,dive,required"`
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
		StatusProperties:  make([]Restrict, 0),
		CommandProperties: make([]*CommandProperty, 0),
	}
}

func (this *DeviceModel) AddModelProperty(restrict Restrict) {
	this.StatusProperties = append(this.StatusProperties, restrict)
}
