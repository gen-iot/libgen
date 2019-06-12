package libgen

import (
	"encoding/json"
	"time"
)

type TimeoutType int

const (
	Timeout TimeoutType = iota
	Custom
	Never
)

const (
	DefaultTimeoutPolicyTimeout = time.Second * 60
)

type RestrictType int

const (
	I32Range RestrictType = iota
	I32Enum
	StrEnum
	Any
)

type PropertyType string

const (
	Status  PropertyType = "status"
	Command PropertyType = "command"
)

type JsonObject map[string]interface{}

func NewJsonObjectFromBytes(data []byte) (JsonObject, error) {
	out := make(JsonObject)
	err := json.Unmarshal(data, &out)
	return out, err
}

func (this JsonObject) GetIntOr(key string, dft int) int {
	i, ok := this[key]
	if !ok {
		return dft
	}
	i2, ok := i.(int)
	if !ok {
		return dft
	}
	return i2
}
