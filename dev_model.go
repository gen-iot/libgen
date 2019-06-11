package libgen

type PropertyType string

const (
	Status  PropertyType = "status"
	Command PropertyType = "command"
)

type LimiterType string

const (
	Any        LimiterType = "any"
	StringEnum LimiterType = "string_enum"
	StringAny  LimiterType = "string_any"
	Int32Enum  LimiterType = "int32_enum"
	Int32Range LimiterType = "int32_range"
)

type PropertyLimiter struct {
	Type        PropertyType
	Name        string
	LimiterType LimiterType
}

type DeviceModel struct {
	Id               string
	Name             string
	PropertyLimiters []PropertyLimiter
}

//type Property struct {
//	LimiterConfig map[string]interface{}
//}

type Device struct {
	ModelId    string
	Id         string
	Name       string
	Room       string
	Properties map[string]interface{}
	MetaData   map[string]interface{}
}
