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

// example modelInfo:
// DeviceModelInfo {
//		Package :"com.example.example",
//		Name	:"light",
//}
type DeviceModelInfo struct {
	Package string `json:"package" validate:"required"` // if empty , use current appid
	Name    string `json:"name" validate:"required"`    // model name
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
