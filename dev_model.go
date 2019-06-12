package libgen

type ModelProperty struct {
	Type    PropertyType
	Name    string
	Restrict Restrict
}

type DeviceModel struct {
	Id              string
	Name            string
	ModelProperties []*ModelProperty
}

type Device struct {
	ModelId    string
	Id         string
	Name       string
	Room       string
	Properties map[string]interface{}
	MetaData   map[string]interface{}
}
