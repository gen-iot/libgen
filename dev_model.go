package libgen

type ModelProperty struct {
	Type     PropertyType `json:"type" validate:"required"`
	Name     string       `json:"name" validate:"required"`
	Restrict Restrict     `json:"restrict" validate:"required"`
}

type DeviceModel struct {
	DeviceModelInfo
	ModelProperties []*ModelProperty `json:"modelProperties"`
}

type DeviceModelInfo struct {
	AppId string `json:"appId"` // if empty , use current appid
	Id    string `json:"modelId" validate:"required"`
	Name  string `json:"name" validate:"required"`
}

type Device struct {
	ModelInfo  *DeviceModelInfo       `json:"modelInfo" validate:"required"`
	Id         string                 `json:"devId" validate:"required"`
	Name       string                 `json:"name" validate:"required"`
	Room       string                 `json:"room" validate:"required"`
	Properties map[string]interface{} `json:"properties"`
	MetaData   map[string]interface{} `json:"metadata"`
}

type DeviceInfo struct {
	Device
	AppId string `json:"appId"`
}
