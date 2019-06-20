package libgen

type ModelProperty struct {
	Type     PropertyType `json:"type" validate:"required"`
	Name     string       `json:"name" validate:"required"`
	Restrict Restrict     `json:"restrict" validate:"required"`
}

type DeviceModel struct {
	Id              string           `json:"modelId" validate:"required"`
	Name            string           `json:"name" validate:"required"`
	ModelProperties []*ModelProperty `json:"modelProperties"`
}

type Device struct {
	ModelId    string                 `json:"modelId" validate:"required"`
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
