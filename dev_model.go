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

type Property struct {
	Type          PropertyType
	Name          string
	LimiterType   LimiterType
	LimiterConfig map[string]interface{}
}

type Device struct {
	Name       string
	Room       string
	Properties []Property
	MetaData   map[string]interface{}
}
